package main

import (
	"context"
	"log"

	"github.com/gerasimovvladislav/telegram"
	"github.com/gerasimovvladislav/telegram/examples/tgbot"
	"github.com/gerasimovvladislav/telegram/examples/tgbot/internal/slash"
	"github.com/gerasimovvladislav/telegram/examples/tgbot/internal/slash/processors"
	"github.com/jessevdk/go-flags"
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
	if err != nil {
		log.Fatal("failed to run bot: ", err)
	}

	err = b.Start(context.Background())
	if err != nil {
		log.Fatal("failed to start bot: ", err)
	}
}
