package framework

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/bwmarrin/discordgo"

	"github.com/endeveit/enca"
	"github.com/endeveit/guesslanguage"

	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

type Bot struct {
	Discord *discordgo.Session

	serverId string
	Routes   []Route

	lg *log.Entry
	db *sql.DB
}

func NewBot(token string, serverId string, db *sql.DB) (*Bot, error) {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	discord.AddHandler(translate)

	return &Bot{
		Discord: discord,

		serverId: serverId,
		Routes:   make([]Route, 0),
		lg:       log.WithField("src", "bot"),
		db:       db,
	}, nil
}

// Register adds routes to the bot
func (b *Bot) Register(routes ...Route) {
	b.Routes = append(b.Routes, routes...)
}

func keyBuilder(opt *discordgo.ApplicationCommandInteractionDataOption) string {
	routeKey := opt.Name

	// Base case: If there are no options of type SubCommand, return the route key
	if len(opt.Options) == 0 {
		return routeKey
	}
	if opt.Options[0].Type != discordgo.ApplicationCommandOptionSubCommand {
		return routeKey
	}

	// Recursive case: If there are options of type SubCommand, append the option name to the route key
	return fmt.Sprintf("%s.%s", routeKey, keyBuilder(opt.Options[0]))
}

func (b *Bot) registerAllCommandsAndRouting() {
	// Register the route with Discord
	for _, route := range b.Routes {
		createdApp, err := b.Discord.ApplicationCommandCreate(b.Discord.State.User.ID, b.serverId, route.declaration)
		if err != nil {
			b.lg.Errorf("Error creating command: %s", err)
			continue
		}
		route.declaration = createdApp

		for k := range route.commandRoute {
			b.lg.Infof("Registered command route: %s", k)
		}
	}

	// Define a function to build the route key
	appKeyBuilder := func(i *discordgo.InteractionCreate) string {
		routeKey := i.ApplicationCommandData().Name

		// Base case: If there are no options of type SubCommand, return the route key
		if len(i.ApplicationCommandData().Options) == 0 {
			return routeKey
		}

		// Recursive case: If there are options of type SubCommand, append the option name to the route key
		if i.ApplicationCommandData().Options[0].Type != discordgo.ApplicationCommandOptionSubCommand {
			return routeKey
		}

		return fmt.Sprintf("%s.%s", routeKey, keyBuilder(i.ApplicationCommandData().Options[0]))
	}

	// Handle the route execution
	b.Discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// TODO: Handle events such as Button interactions with Custom IDs etc.
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		// Build the Route Key
		routeKey := appKeyBuilder(i)

		// Find the route
		for _, route := range b.Routes {
			if er, ok := route.commandRoute[routeKey]; ok {
				b.lg.Infof("Executing route: %s", routeKey)
				er.Execute(NewContext(
					WithSession(s),
					WithInteraction(i.Interaction),
					WithMessage(i.Interaction.Message),
					WithLogger(b.lg.WithField("route", routeKey)),
					WithDatabase(b.db),
				))
				return
			}
		}

		// If the route is not found, respond with an error message
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Command not found",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	})
}

func (b *Bot) deregisterAllCommands() {

	// Delete the route from Discord
	for _, route := range b.Routes {
		b.Discord.ApplicationCommandDelete(b.Discord.State.User.ID, b.serverId, route.declaration.ID)

		for k := range route.commandRoute {
			b.lg.Infof("Deregistered command route: %s", k)
		}
	}
}

func (b *Bot) DefineModerationRules(rules ...ActionableRule) {
	b.Discord.AddHandler(
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
						rule.Rule.Action(NewContext(
							WithSession(s),
							WithInteraction(nil), // No interaction for messages
							WithMessage(m.Message),
							WithLogger(b.lg.WithField("rule", rule.Rule.Name())),
							WithDatabase(b.db),
						), err)
					}
				}
			}
		})
}

func (b *Bot) Run() error {
	if err := b.Discord.Open(); err != nil {
		return err
	}
	b.registerAllCommandsAndRouting()
	return nil
}

func (b *Bot) Close() error {
	b.deregisterAllCommands()
	return b.Discord.Close()
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
