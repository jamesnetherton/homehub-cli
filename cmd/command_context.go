package cmd

import "strconv"
import "github.com/chzyer/readline"

// CommandContext encapsulates the state of a command execution lifecycle
type CommandContext struct {
	args   []string
	result interface{}
	err    error
}

// IsError returns whether any step within the command execution lifecycle encountered an error
func (c *CommandContext) IsError() bool {
	return c.err != nil
}

// HasResult returns whether any step within the command execution returned a valid response
func (c *CommandContext) HasResult() bool {
	return c.GetResult() != nil
}

// GetIntArg returns an int cli command argument at the specified index
func (c *CommandContext) GetIntArg(index int) (result int, err error) {
	return strconv.Atoi(c.args[index])
}

// GetStringArg returns a string cli command argument at the specified index
func (c *CommandContext) GetStringArg(index int) (result string) {
	return c.args[index]
}

// GetBooleanArg returns a bool cli command argument at the specified index
func (c *CommandContext) GetBooleanArg(index int) (result bool, err error) {
	return strconv.ParseBool(c.args[index])
}

// GetResult returns the result of the command execution lifecycle
func (c *CommandContext) GetResult() interface{} {
	return c.result
}

// SetResult sets the result of the command execution lifecycle
func (c *CommandContext) SetResult(result interface{}, err error) {
	c.result = result
	c.err = err
}

// ReadLine reads a line from the command CLI
func (c *CommandContext) ReadLine(prompt string) (result string, err error) {
	readLine, err := readline.New(prompt)
	if err == nil {
		defer readLine.Close()
		return readLine.Readline()
	}
	return "", err
}

// ReadPassword reads a password string from the command CLI
func (c *CommandContext) ReadPassword(prompt string) (result string, err error) {
	readLine, err := readline.New(prompt)
	if err == nil {
		defer readLine.Close()
		password, err := readLine.ReadPassword(prompt)
		return string(password[:]), err
	}
	return "", err
}
