package main

type Config struct {
	BotToken string `long:"bot-token" description:"Bot token" env:"BOT_TOKEN" required:"true"`
	AdminID  int64  `long:"admin-id" description:"Admin id" env:"ADMIN_ID" required:"true"`
	UserID   int64  `long:"user-id" description:"User id" env:"USER_ID" required:"true"`
}
