package cmd

import "strconv"

type CommandContext struct {
	args   []string
	result interface{}
	err    error
}

func (c *CommandContext) IsError() bool {
	return c.err != nil
}

func (c *CommandContext) HasResult() bool {
	return c.GetResult() != nil
}

func (c *CommandContext) GetIntArg(index int) (result int, err error) {
	return strconv.Atoi(c.args[index])
}

func (c *CommandContext) GetStringArg(index int) (result string) {
	return c.args[index]
}

func (c *CommandContext) GetBooleanArg(index int) (result bool, err error) {
	return strconv.ParseBool(c.args[index])
}

func (c *CommandContext) GetResult() interface{} {
	return c.result
}

func (c *CommandContext) SetResult(result interface{}, err error) {
	c.result = result
	c.err = err
}
