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
	PreExec     func(args []string) error
	Exec        func(args []string) (result interface{}, err error)
	PostExec    func(result interface{}, err error) error
}

// ExecuteLifecylce runs the command execution lifecycle
func (c *GenericCommand) ExecuteLifecylce(args []string) {

	if helpRequested(args) {
		c.Explain()
		return
	}

	if c.Validate(args) {
		if c.PreExec == nil {
			c.PreExec = func(args []string) error {
				return nil
			}
		}

		if c.PostExec == nil {
			c.PostExec = func(result interface{}, err error) error {
				if err != nil {
					return err
				} else {
					if result != nil {
						fmt.Println(result)
					}
					return nil
				}
			}
		}

		preErr := c.PreExec(args)

		if preErr == nil {
			postErr := c.PostExec(c.Execute(args))
			if postErr != nil {
				fmt.Println(postErr)
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
func (c *GenericCommand) Validate(args []string) bool {
	if len(args) != len(c.ArgNames) {
		return false
	}
	return true
}

// Execute executes the command
func (c *GenericCommand) Execute(args []string) (result interface{}, err error) {
	return c.Exec(args)
}

func (c *GenericCommand) Usage() {
	fmt.Printf("Usage: %s ", c.Name)

	for i := 0; i < len(c.ArgNames); i++ {
		fmt.Printf("%s<%s> ", c.ArgNames[i], c.ArgTypes[i])
	}

	fmt.Println()
}

func (c *GenericCommand) Explain() {
	fmt.Printf("%s: %s\n\n", c.Name, c.Description)
	c.Usage()
}

func helpRequested(args []string) bool {
	return (len(args) > 0) && (args[0] == "-help" || args[0] == "--help" || args[0] == "-h")
}
