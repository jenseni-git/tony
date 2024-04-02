package framework

type ModerateRule interface {
	// Name of the rule
	Name() string

	// Test the rule against content
	Test(content string) error

	// What action should be taken if the rule is violated
	Action(ctx *Context, violation error)
}

type ActionableRule struct {
	Channel string
	Rule    ModerateRule
}

func Rule(channel string, rule ModerateRule) ActionableRule {
	return ActionableRule{
		Channel: channel,
		Rule:    rule,
	}
}
