package commands

import (
	"github.com/aussiebroadwan/tony/framework"
	"github.com/bwmarrin/discordgo"
)

// PingCommand defines a command for responding to "ping" interactions
// with a simple "Pong!" message. This command demonstrates a basic
// interaction within Discord using the discordgo package.
type PingCommand struct {
	framework.Command
}

// Register is responsible for registering the "ping" command with
// Discord's API. It defines the command name and description that
// appear in the Discord user interface.
func (pc *PingCommand) Register(s *discordgo.Session) *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Replies with Pong!",
	}
}

// Execute handles the execution logic for the "ping" command. When a user
// invokes this command, Discord triggers this method, allowing the bot to
// respond appropriately.
func (pc *PingCommand) Execute(ctx *framework.Context) {
	ctx.Session().InteractionRespond(ctx.Interaction(), &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})
}

func (pc *PingCommand) OnEvent(ctx *framework.Context, eventType framework.EventType) {
	// NOP
}