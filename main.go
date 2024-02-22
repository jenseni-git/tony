package main

import (
	"os"
	"os/signal"

	"github.com/aussiebroadwan/tony/framework"
	"github.com/aussiebroadwan/tony/internal/moderation"

	log "github.com/sirupsen/logrus"
)

const VERSION = "0.1.0"

var (
	token    = os.Getenv("DISCORD_TOKEN")
	serverId = os.Getenv("DISCORD_SERVER_ID")
)

func init() {
	// Setup logging
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	// Check if token is provided
	if token == "" {
		log.Fatal("No token provided. Please set DISCORD_TOKEN environment variable.")
		return
	}
}

func main() {
	// Create a new bot
	bot, err := framework.NewBot(token, serverId)
	if err != nil {
		log.Fatalf("Error creating bot: %s", err)
		return
	}

	bot.DefineModerationRules(
		framework.Rule("tech-news", &moderation.ModerateNewsRule{}),
		framework.Rule("rss", &moderation.ModerateRSSRule{}),
	)

	// Run the bot
	if err = bot.Run(); err != nil {
		log.Fatalf("Error running bot: %s", err)
		return
	}
	defer bot.Close()

	waitForInterrupt()
}

func waitForInterrupt() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Shutting down...")
}
