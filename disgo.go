package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("Invalid args.\nUsage:\n./main <path/to/config>")
		log.Fatal("Exiting...")
	}

	// load bot configs from cmdline file
	c, err := loadBotConfig(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to parse configs: %v", err)
	}

	// create bot session with token
	s, err := discordgo.New("Bot " + c.Token)
	if err != nil {
		log.Fatalf("Couldn't create session: %v", err)
	}

	// add via feature flags
	s.AddHandler(LuckyWheelHandler)
	s.AddHandler(ready)

	// 'Open()' starts the bot session
	if err = s.Open(); err != nil {
		log.Fatalf("Couldn't open session: %v", err)
	}
	defer s.Close()

	// add intents via command requirement
	//s.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildPresences
	s.Identify.Intents = discordgo.IntentsAll

	// bot blocks until kill signal recieved
	kill := make(chan os.Signal)
	signal.Notify(kill, os.Interrupt)
	<-kill
	log.Println("Shuting down...")
}

func ready(s *discordgo.Session, r *discordgo.Ready) {
	log.Println("Bot is up!")
}
