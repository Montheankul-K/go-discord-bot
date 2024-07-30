package server

import (
	"flag"
	"github.com/Montheankul-K/go-discord-bot/configs"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
)

var (
	GuildID        = flag.String("guide", "", "test guild ID. if not passed - bot registers commands globally")
	RemoveCommands = flag.Bool("rmcmd", true, "remove all commands after shutdown or not")
)

type IDiscordServer interface {
	Start()
}

type discordServer struct {
	cfg      configs.IConfig
	discord  *discordgo.Session
	commands []*discordgo.ApplicationCommand
}

func NewDiscordServer(cfg configs.IConfig) IDiscordServer {
	discord, err := discordgo.New("Bot " + cfg.App().GetToken())
	if err != nil {
		log.Fatal("failed to create discord session:", err)
	}

	return &discordServer{
		cfg:      cfg,
		discord:  discord,
		commands: make([]*discordgo.ApplicationCommand, 0),
	}
}

func (s *discordServer) Start() {
	s.discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("logged in as %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	if err := s.discord.Open(); err != nil {
		log.Fatalf("failed to open discord session: %s", err.Error())
	}

	m := ModuleInit(s)
	m.BotinfoModule().Init()

	registeredCommands := make([]*discordgo.ApplicationCommand, len(s.commands))
	for i, v := range s.commands {
		cmd, err := s.discord.ApplicationCommandCreate(s.discord.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("failed to create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	defer func(discord *discordgo.Session) {
		err := discord.Close()
		if err != nil {
			log.Fatalf("failed to close discord session: %s", err.Error())
		}
	}(s.discord)

	s.discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := m.GetCommandHandlers()[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("press ctrl + c to exit")
	<-stop

	if *RemoveCommands {
		log.Println("removing commands")
		for _, v := range registeredCommands {
			err := s.discord.ApplicationCommandDelete(s.discord.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("failed to delete '%v' command: %v", v.Name, err)
			}
		}
	}
	log.Println("gracefully shutting down")
}
