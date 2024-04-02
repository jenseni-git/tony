package commands

import (
	"fmt"
	"time"

	"github.com/aussiebroadwan/tony/database"
	"github.com/aussiebroadwan/tony/framework"
	"github.com/aussiebroadwan/tony/pkg/reminders"
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
	/* NOP */
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
	interaction := ctx.Interaction()
	db := ctx.Database()
	commandOptions := interaction.ApplicationCommandData().Options[0].Options

	// Get the time and message from the interaction
	triggerTimeStr, err := framework.GetOption(commandOptions, "time")
	if err != nil {
		ctx.Session().InteractionRespond(ctx.Interaction(), &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Error: Time is required",
			},
		})
		return
	}

	message, err := framework.GetOption(commandOptions, "message")
	if err != nil {
		ctx.Session().InteractionRespond(ctx.Interaction(), &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Error: Message is required",
			},
		})
		return
	}

	// Check if the time is valid
	triggerTime, err := time.Parse(time.DateTime, triggerTimeStr.StringValue())
	if err != nil {
		ctx.Session().InteractionRespond(ctx.Interaction(), &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Error: Invalid time format eg. 2022-01-01 15:04:05",
			},
		})
		return
	}

	// Add the reminder
	err = database.AddReminder(
		db,
		interaction.User.Username,
		triggerTime,
		ctx.Session(),
		interaction.ChannelID,
		message.StringValue(),
	)

	if err != nil {
		ctx.Session().InteractionRespond(ctx.Interaction(), &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Error: Failed to add reminder",
			},
		})
		return
	}
}

func (com *RemindAddSubCommand) OnEvent(ctx *framework.Context, eventType framework.EventType) {
	/* NOP */
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
	interaction := ctx.Interaction()
	db := ctx.Database()
	commandOptions := interaction.ApplicationCommandData().Options[0].Options

	// Get the reminder ID from the interaction
	id, err := framework.GetOption(commandOptions, "id")
	if err != nil {
		ctx.Session().InteractionRespond(ctx.Interaction(), &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Error: id is required",
			},
		})
		return
	}

	// Delete the reminder
	err = database.DeleteReminder(db, id.IntValue())
	if err != nil {
		ctx.Session().InteractionRespond(ctx.Interaction(), &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf("Error: reminder %d not found", id.IntValue()),
			},
		})
	}
}

func (c *RemindDeleteSubCommand) OnEvent(ctx *framework.Context, eventType framework.EventType) {
	/* NOP */
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
	interaction := ctx.Interaction()

	// Get all reminders
	reminderList := reminders.List()

	// Fetch current user's reminders
	var userReminders = make([]reminders.Reminder, 0)
	for _, reminder := range reminderList {
		if reminder.CreatedBy == interaction.User.Username {
			userReminders = append(userReminders, reminder)
		}
	}

	// Respond with the reminder list
	if len(userReminders) == 0 {
		ctx.Session().InteractionRespond(ctx.Interaction(), &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "No reminders found",
			},
		})
		return
	}

	var reminderListStr string = "Reminders:\n\n```\n"
	for _, reminder := range userReminders {
		reminderListStr += fmt.Sprintf("ID: %d, Time: %s\n", reminder.ID, reminder.TriggerTime.String())
	}
	reminderListStr += "```"

	ctx.Session().InteractionRespond(ctx.Interaction(), &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: reminderListStr,
		},
	})
}

func (c *RemindListSubCommand) OnEvent(ctx *framework.Context, eventType framework.EventType) {
	/* NOP */
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
	interaction := ctx.Interaction()
	commandOptions := interaction.ApplicationCommandData().Options[0].Options

	// Get the reminder ID from the interaction
	id, err := framework.GetOption(commandOptions, "id")
	if err != nil {
		ctx.Session().InteractionRespond(ctx.Interaction(), &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Error: id is required",
			},
		})
		return
	}

	// Get the reminder status
	timeLeft, err := reminders.Status(id.IntValue())
	if err != nil {
		ctx.Session().InteractionRespond(ctx.Interaction(), &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf("Error: reminder %d not found", id.IntValue()),
			},
		})
		return
	}

	// Respond with the reminder status
	ctx.Session().InteractionRespond(ctx.Interaction(), &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf("Time left: %s", timeLeft.String()),
		},
	})
}

func (c *RemindStatusSubCommand) OnEvent(ctx *framework.Context, eventType framework.EventType) {
	/* NOP */
}
