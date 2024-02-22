package framework

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
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

	return &Bot{
		discord:  discord,
		serverId: serverId,
	}, nil
}

func (b *Bot) DefineModerationRules(rules ...ActionableRule) {
	b.discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	})
}

func (b *Bot) Run() error {
	return b.discord.Open()
}

func (b *Bot) Close() error {
	return b.discord.Close()
}
