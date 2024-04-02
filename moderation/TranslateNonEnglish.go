package moderation

import (
	"fmt"

	"github.com/aussiebroadwan/tony/framework"
	"github.com/endeveit/enca"
	"github.com/endeveit/guesslanguage"
)

type TranslateNonEnglish struct {
	framework.ModerateRule
}

func (r *TranslateNonEnglish) Name() string {
	return "translate-non-english"
}

// Test tests the rule against the content
func (r *TranslateNonEnglish) Test(ctx *framework.Context, content string) error {
	lg := ctx.Logger().WithField("moderation", "translate")
	lg.Infof("Testing content: %s", content)

	// Check if the content is in English or not
	lang, err := guesslanguage.Guess(content)
	if err == nil && !(lang == "en" || lang == "UNKNOWN") {
		return fmt.Errorf(lang)
	}

	// Log the error if there is one
	if err != nil {
		ctx.Logger().WithError(err).Error("Error guessing language")
	}

	return nil
}

// Action takes action if the rule is violated, in this case, it translates the
// non-English text
func (r *TranslateNonEnglish) Action(ctx *framework.Context, violation error) {
	lg := ctx.Logger().WithField("moderation", "translate")

	// Get the language of the text from the original test
	language := violation.Error()

	// Create a new enca analyzer
	analyzer, err := enca.New(language)
	if err != nil {
		lg.WithError(err).Error("Error creating enca analyzer")
		return
	}

	// Analyze the text
	encoding, err := analyzer.FromString(ctx.Message().Content, enca.NAME_STYLE_HUMAN)
	defer analyzer.Free()

	// And with no errors, print out the translation
	if err != nil {
		lg.WithError(err).Error("Error analyzing text")
		return
	}

	// Send the message to the channel
	outMessage := fmt.Sprintf("This message likely contained %s text\nIts English meaning is: %s", language, encoding)
	ctx.Session().ChannelMessageSend(ctx.Message().ChannelID, outMessage)
}
