package cmd

import "fmt"

var date string
var firmware string
var revision string
var version string

// NewAboutCommand creates a new command to output information about the homehub-cli
func NewAboutCommand() *GenericCommand {
	return &GenericCommand{
		Name:        "About",
		Description: "Displays information about the homehub-cli",
		Exec: func(context *CommandContext) {
			fmt.Printf("%-20s%s\n", "Version:", version)
			fmt.Printf("%-20s%s\n", "Compatible with:", firmware)
			fmt.Printf("%-20s%s\n", "Build date:", date)
			fmt.Printf("%-20s%s\n", "Revision:", revision)
		},
	}
}
