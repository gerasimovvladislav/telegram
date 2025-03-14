package main

import (
	"context"
	"log"

	"github.com/jessevdk/go-flags"
	"gitlab.com/vladislavgerasimov/telegram"
	"gitlab.com/vladislavgerasimov/telegram/examples/tgbot"
	"gitlab.com/vladislavgerasimov/telegram/examples/tgbot/internal/slash"
	"gitlab.com/vladislavgerasimov/telegram/examples/tgbot/internal/slash/processors"
)

func main() {
	var cfg Config
	f := flags.NewParser(&cfg, flags.Default)
	_, err := f.Parse()
	if err != nil {
		log.Fatal("failed to parse config: ", err)
	}

	var b telegram.Bot
	b, err = tgbot.Init(
		cfg.BotToken,
		[]int64{
			cfg.AdminID,
			cfg.UserID,
		},
		nil,
		map[telegram.SlashCommandId]telegram.SlashCommandProcessor{
			slash.CommandIdHello: processors.NewHello(),
		},
		nil,
		nil,
	)
	err = b.Start(context.Background())
	if err != nil {
		log.Fatal("failed to run bot: ", err)
	}
}
