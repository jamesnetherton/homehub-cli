package cmd

// Command is an interface that defines a Hub CLI Command
type Command interface {
	Execute(context *CommandContext)
	ExecuteLifecylce(args []string)
	GetName() string
	Usage()
	Validate(context *CommandContext) bool
}
