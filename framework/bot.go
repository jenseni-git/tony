package framework

import (
	"fmt"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/endeveit/enca"
	"github.com/endeveit/guesslanguage"
)

type Bot struct {
	discord  *discordgo.Session
	serverId string
}

func NewBot(token string, serverId string) (*Bot, error) {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	discord.AddHandler(translate)

	return &Bot{
		discord:  discord,
		serverId: serverId,
	}, nil
}

func (b *Bot) DefineModerationRules(rules ...ActionableRule) {
	b.discord.AddHandler(
		func(s *discordgo.Session, m *discordgo.MessageCreate) {
			if m.Author.ID == s.State.User.ID {
				return
			}

			// Get Channel Name from Message
			channel, err := s.State.Channel(m.ChannelID)
			if err != nil {
				return
			}

			// Test a regex match for the channel name against the rule
			for _, rule := range rules {
				if match, _ := regexp.Match(rule.Channel, []byte(channel.Name)); match {

					// Test the rule
					if err := rule.Rule.Test(m.Content); err != nil {
						rule.Rule.Action(&Context{
							Message: m.Message,
							Session: s,
						}, err)
					}
				}
			}
		},
	)
}

func (b *Bot) Run() error {
	return b.discord.Open()
}

func (b *Bot) Close() error {
	return b.discord.Close()
}

func translate(session *discordgo.Session, message *discordgo.MessageCreate) {
	fmt.Println("translate() received message:\"" + message.Content + "\"")

	// Check all non-bot messages
	if message.Author.ID == session.State.User.ID {
		return
	}
	// Check for non-english messages first to prevent unneccessary translations
	lang, err := guesslanguage.Guess(message.Content)
	if err != nil && lang != "en" {
		// If non-english, translate
		analyzer, err := enca.New(lang)

		if err == nil {
			encoding, err := analyzer.FromString(message.Content, enca.NAME_STYLE_HUMAN)
			defer analyzer.Free()

			// And with no errors, print out the translation
			if err == nil {
				out_message := "This message likely contained " + lang +
					" text\nIts English meaning is: " + encoding
				session.ChannelMessageSend(message.ChannelID, out_message)
			}
		}
	}
	return
}
