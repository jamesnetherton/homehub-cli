package cmd

import (
	"fmt"
)

// GenericCommand defines a Hub CLI command that does not require authentication prior to execution
type GenericCommand struct {
	Name        string
	Description string
	ArgNames    []string
	ArgTypes    []string
	PreExec     func(context *CommandContext)
	Exec        func(context *CommandContext)
	PostExec    func(context *CommandContext)
}

// ExecuteLifecylce runs the command execution lifecycle
func (c *GenericCommand) ExecuteLifecylce(args []string) {

	context := &CommandContext{
		args: args,
	}

	if helpRequested(context) {
		c.Explain()
		return
	}

	if c.Validate(context) {
		// Default PreExec
		if c.PreExec == nil {
			c.PreExec = func(context *CommandContext) {}
		}

		// Default PostExec
		if c.PostExec == nil {
			c.PostExec = func(context *CommandContext) {
				if !context.IsError() && context.HasResult() {
					fmt.Println(context.GetResult())
				}
			}
		}

		c.PreExec(context)

		if !context.IsError() {
			c.Execute(context)
			c.PostExec(context)
			if context.IsError() {
				fmt.Println(context.err)
			}
		}
	} else {
		c.Usage()
	}
}

// GetName returns the name of the command
func (c *GenericCommand) GetName() string {
	return c.Name
}

// Validate validates that the correct number of arguments were passed to the command
func (c *GenericCommand) Validate(context *CommandContext) bool {
	if len(context.args) != len(c.ArgNames) {
		return false
	}
	return true
}

// Execute executes the command
func (c *GenericCommand) Execute(context *CommandContext) {
	c.Exec(context)
}

func (c *GenericCommand) Usage() {
	fmt.Printf("Usage: %s ", c.Name)

	for i := 0; i < len(c.ArgNames); i++ {
		fmt.Printf("%s <%s> ", c.ArgNames[i], c.ArgTypes[i])
	}

	fmt.Println()
}

func (c *GenericCommand) Explain() {
	fmt.Printf("%s: %s\n\n", c.Name, c.Description)
	c.Usage()
}

func helpRequested(context *CommandContext) bool {
	return (len(context.args) > 0) && (context.GetStringArg(0) == "-help" || context.GetStringArg(0) == "--help" || context.GetStringArg(0) == "-h")
}
