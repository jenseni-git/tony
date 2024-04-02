package framework

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// EventType is a custom type to define various event types a command can handle
type EventType string

const (
	CommandType EventType = "command"
	ButtonType  EventType = "button"
	SelectType  EventType = "select"
	// Additional event types can be added here
)

type Command interface {
	Execute(ctx *Context)                      // Executed for slash commands
	OnEvent(ctx *Context, eventType EventType) // Handles various event types
}

// Command interface now includes OnEvent instead of OnButton and OnSelect
type AppCommand interface {
	Register(s *discordgo.Session) *discordgo.ApplicationCommand
	Command
}

type SubCommand interface {
	Register(s *discordgo.Session) *discordgo.ApplicationCommandOption
	Command
}

// Route associates a command name with a command instance and optional subcommands
type Route struct {
	Name      string
	Executes  bool
	Command   AppCommand
	SubRoutes []SubRoute

	declaration  *discordgo.ApplicationCommand
	commandRoute map[string]Command
}

// NewRoute constructs a new Route
func NewRoute(bot *Bot, name string, executionLogic bool, command AppCommand, subroutes ...SubRoute) Route {
	r := Route{
		Name:         name,
		Executes:     executionLogic,
		Command:      command,
		SubRoutes:    subroutes,
		declaration:  command.Register(bot.discord),
		commandRoute: make(map[string]Command),
	}

	for _, sr := range subroutes {
		for k, v := range sr.commandRoute {
			r.commandRoute[fmt.Sprintf("%s.%s", r.Name, k)] = v
		}

		r.declaration.Options = append(r.declaration.Options, sr.declaration)
	}

	if executionLogic {
		r.commandRoute[name] = command
	}

	return r
}

type SubRoute struct {
	Name       string
	Executes   bool
	SubCommand SubCommand
	SubRoutes  []SubRoute

	declaration  *discordgo.ApplicationCommandOption
	commandRoute map[string]Command
}

func NewSubRoute(bot *Bot, name string, executionLogic bool, subcommand SubCommand, subroutes ...SubRoute) SubRoute {
	r := SubRoute{
		Name:       name,
		Executes:   executionLogic,
		SubCommand: subcommand,
		SubRoutes:  subroutes,

		declaration:  subcommand.Register(bot.discord),
		commandRoute: make(map[string]Command),
	}

	// Check if the subroute has subroutes
	for _, sr := range subroutes {

		// Add the subcommand to the command route
		for k, v := range sr.commandRoute {
			r.commandRoute[fmt.Sprintf("%s.%s", r.Name, k)] = v
		}

		// Add the subcommands to the declaration
		r.declaration.Options = append(r.declaration.Options, sr.declaration)
	}

	// Add the subcommand to the command route
	if executionLogic {
		r.commandRoute[name] = subcommand
	}

	return r
}