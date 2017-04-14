package cmd

// Command is an interface that defines a Hub CLI Command
type Command interface {
	Execute(args []string) (result interface{}, err error)
	ExecuteLifecylce(args []string)
	GetName() string
	Usage()
	Validate(args []string) bool
}
