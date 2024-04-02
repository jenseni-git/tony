package moderation

import (
	"errors"
	"regexp"

	"github.com/aussiebroadwan/tony/framework"
)

// The #rss channel is for news from RSS feeds links. We like to keep the channel
// clean and only show RSS links with a title and description. This rule will
// moderate the news in that channel. More specifically it will check and enforce
// the rss post format.
//
// The format is as follows:
//
//	**<title>**: <link>
//
//	{description}
//
// Example:
//
//	**Smashing Magazine**: http://smashingmagazine.com/feed
//
//	This one is mostly for the frontend dev, design, and UX nerds. Smashing Mag
//	has some really good content about design in general, as well as plenty of
//	tips and tricks for improving your workflow and usability of your work.
//
// If the message does not match the format, the bot will delete the message and
// send a message to the user to let them know that the message was deleted and
// why.
type ModerateRSSRule struct {
	framework.ModerateRule
}

var (
	ErrInvalidRSSPostFormat = errors.New("rss posts must be in the following format:\n```\n**<title>**: <link>\n\n{description}```")
	ErrRSSTitleFormatError  = errors.New("the title must be in bold and end with a colon then a link")
)

func (r *ModerateRSSRule) Name() string {
	return "rss"
}

// Test tests the rule against the content
func (r *ModerateRSSRule) Test(ctx *framework.Context, content string) error {
	// Check if the title is in bold and ends with a colon and a link
	blockRegex := regexp.MustCompile(`\*\*.*\*\*: http(s)?:\/\/.*\n\n.*`)
	if !blockRegex.MatchString(content) {
		return ErrInvalidRSSPostFormat
	}

	return nil
}

// Action takes action if the rule is violated
func (r *ModerateRSSRule) Action(ctx *framework.Context, violation error) {
	// Delete the message
	ctx.Session().ChannelMessageDelete(ctx.Message().ChannelID, ctx.Message().ID)

	// Get or create a DM channel with the user
	dmChannel, err := ctx.Session().UserChannelCreate(ctx.Message().Author.ID)
	if err != nil {
		// Handle error, log it, or take appropriate action
		return
	}

	// Send a direct message to the user
	ctx.Session().ChannelMessageSend(dmChannel.ID, violation.Error())
}
