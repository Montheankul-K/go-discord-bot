package main

import (
	"github.com/Montheankul-K/go-discord-bot/configs"
	"github.com/Montheankul-K/go-discord-bot/modules/server"
)

func main() {
	cfg := configs.NewConfig("./.env")
	server.NewDiscordServer(cfg).Start()
}
