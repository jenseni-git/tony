package framework

import "github.com/bwmarrin/discordgo"

type Context struct {
	// The message that triggered the command
	Message *discordgo.Message

	// The discord session
	Session *discordgo.Session
}
