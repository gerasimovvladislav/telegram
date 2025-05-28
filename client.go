package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type sendJob struct {
	msg    tgbotapi.Chattable
	respCh chan response
}

type response struct {
	msg tgbotapi.Message
	err error
}

type Client struct {
	*tgbotapi.BotAPI
	transport *http.Transport

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	sendCh chan sendJob
}

// runSender sends messages
func (c *Client) runSender() {
	defer c.wg.Done()

	sleepUntil := time.Now()

	for {
		select {
		case <-c.ctx.Done():
			return
		case job := <-c.sendCh:
			if time.Now().Before(sleepUntil) {
				wait := time.Until(sleepUntil)
				time.Sleep(wait)
			}

			msg, err := c.BotAPI.Send(job.msg)
			if err != nil {
				if IsFloodError(err) {
					retryAfter := ParseRetryAfter(err)
					if retryAfter <= 0 {
						retryAfter = 3
					}
					sleepUntil = time.Now().Add(time.Duration(retryAfter) * time.Second)
					c.sendCh <- job
					continue
				}

				job.respCh <- response{err: err}
				continue
			}

			job.respCh <- response{msg: msg, err: nil}

			time.Sleep(500 * time.Millisecond)
		}
	}
}

// NewClient creates new client
func NewClient(token string) (*Client, error) {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 5 * time.Second,
		MaxIdleConns:        100,
		IdleConnTimeout:     30 * time.Second,
		ForceAttemptHTTP2:   true,
	}

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	bot, err := tgbotapi.NewBotAPIWithClient(token, tgbotapi.APIEndpoint, httpClient)
	if err != nil {
		return nil, err
	}
	bot.Debug = false

	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		BotAPI:    bot,
		transport: transport,
		ctx:       ctx,
		cancel:    cancel,
		sendCh:    make(chan sendJob, 100),
	}

	client.wg.Add(1)
	go client.runSender()

	return client, nil
}

// Send sends message
func (c *Client) Send(cmsg tgbotapi.Chattable) (tgbotapi.Message, error) {
	respCh := make(chan response, 1)
	c.sendCh <- sendJob{msg: cmsg, respCh: respCh}
	resp := <-respCh
	return resp.msg, resp.err
}

func (c *Client) Shutdown() {
	c.cancel()
	c.wg.Wait()
	c.transport.CloseIdleConnections()
}

// Start starts the bot
func (c *Client) Start(ctx context.Context, controller func(update *Update) error) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	var wg sync.WaitGroup

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
		}

		updates, err := c.GetUpdatesWithContext(ctx, u)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				break loop
			}
			if errors.Is(err, context.DeadlineExceeded) {
				continue
			}

			time.Sleep(3 * time.Second)
			continue
		}

		for _, update := range updates {
			if update.UpdateID >= u.Offset {
				u.Offset = update.UpdateID + 1

				wg.Add(1)

				go func(ctx context.Context, update tgbotapi.Update) {
					defer wg.Done()

					select {
					case <-ctx.Done():
						return
					default:
					}

					wrappedUpdate := WrapUpdate(&update)

					if userId := wrappedUpdate.UserID(); userId == 0 {
						return
					}

					if !wrappedUpdate.Processable() {
						return
					}

					defer func() {
						//TODO: handle error
						_ = recover()
					}()

					//TODO: handle error
					_ = controller(wrappedUpdate)
				}(ctx, update)
			}
		}
	}

	wg.Wait()

	c.transport.CloseIdleConnections()

	return nil
}

// GetUpdatesWithContext returns updates
func (c *Client) GetUpdatesWithContext(ctx context.Context, config tgbotapi.UpdateConfig) ([]tgbotapi.Update, error) {
	params := make(map[string]string)
	params["timeout"] = strconv.Itoa(config.Timeout)
	params["offset"] = strconv.Itoa(config.Offset)
	if config.Limit != 0 {
		params["limit"] = strconv.Itoa(config.Limit)
	}
	if config.AllowedUpdates != nil {
		allowedUpdatesBytes, _ := json.Marshal(config.AllowedUpdates)
		params["allowed_updates"] = string(allowedUpdatesBytes)
	}
	if config.Timeout == 0 {
		params["timeout"] = "60"
	}

	req, err := c.BotAPI.MakeRequest("getUpdates", params)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var updates []tgbotapi.Update
	err = json.Unmarshal(req.Result, &updates)
	if err != nil {
		return nil, err
	}
	return updates, nil
}
