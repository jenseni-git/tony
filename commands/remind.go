package commands

import (
	"github.com/aussiebroadwan/tony/framework"
	"github.com/bwmarrin/discordgo"
)

type RemindCommand struct {
	framework.Command
}

// Register is responsible for registering the "remind" command with
// Discord's API. It defines the command name and description that
// appear in the Discord user interface.
func (c *RemindCommand) Register(s *discordgo.Session) *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "remind",
		Description: "Allows users to set reminders",
	}
}

func (c *RemindCommand) Execute(ctx *framework.Context) {
	// NOP
	//
	// This command is a parent command and does not have any  execution logic.
	// It is used to group subcommands.
}

func (c *RemindCommand) OnEvent(ctx *framework.Context, eventType framework.EventType) {
	// NOP
}

// remind add <time> <message>
type RemindAddSubCommand struct {
	framework.SubCommand
}

func (c *RemindAddSubCommand) Register(s *discordgo.Session) *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "add",
		Description: "Add a reminder",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "message",
				Description: "The message to remind you about",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "time",
				Description: "The time to remind you",
				Required:    true,
			},
		},
	}
}

func (c *RemindAddSubCommand) Execute(ctx *framework.Context) {
	// NOP
}

func (com *RemindAddSubCommand) OnEvent(ctx *framework.Context, eventType framework.EventType) {
	// NOP
}

// remind del <id>
type RemindDeleteSubCommand struct {
	framework.SubCommand
}

func (c *RemindDeleteSubCommand) Register(s *discordgo.Session) *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "del",
		Description: "Delete a reminder",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "The ID of the reminder to delete",
				Required:    true,
			},
		},
	}
}

func (c *RemindDeleteSubCommand) Execute(ctx *framework.Context) {
	// NOP
}

func (c *RemindDeleteSubCommand) OnEvent(ctx *framework.Context, eventType framework.EventType) {
	// NOP
}

// remind list
type RemindListSubCommand struct {
	framework.SubCommand
}

func (c *RemindListSubCommand) Register(s *discordgo.Session) *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "list",
		Description: "List all reminders",
	}
}

func (c *RemindListSubCommand) Execute(ctx *framework.Context) {
	// NOP
}

func (c *RemindListSubCommand) OnEvent(ctx *framework.Context, eventType framework.EventType) {
	// NOP
}

// remind status <id>
type RemindStatusSubCommand struct {
	framework.SubCommand
}

func (c *RemindStatusSubCommand) Register(s *discordgo.Session) *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "status",
		Description: "Get the status of a reminder",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "The ID of the reminder to check",
				Required:    true,
			},
		},
	}
}

func (c *RemindStatusSubCommand) Execute(ctx *framework.Context) {
	// NOP
}

func (c *RemindStatusSubCommand) OnEvent(ctx *framework.Context, eventType framework.EventType) {
	// NOP
}
