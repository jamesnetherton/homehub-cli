package cmd

import (
	"fmt"

	"github.com/bgentry/speakeasy"
	"github.com/jamesnetherton/homehub-cli/service"
)

// NewLoginCommand creates a new command to invoke the Hub Login function
func NewLoginCommand() *GenericCommand {
	return &GenericCommand{
		Name:        "Login",
		Description: "Creates a new Home Hub login session",
		Exec: func(context *CommandContext) {
			hub := service.GetHub()
			if hub == nil || !service.IsLoggedIn() {
				var hubURL string
				var userName string

				fmt.Print("Home hub URL: ")
				fmt.Scan(&hubURL)

				fmt.Print("Home hub user: ")
				fmt.Scan(&userName)

				password, _ := speakeasy.Ask(fmt.Sprint("Home Hub password: "))

				service.NewHub(hubURL, userName, password)
			}
			success, err := service.GetHub().Login()
			if success {
				service.AuthenticationComplete()
			}

			context.SetResult(success, err)
		},
		PostExec: func(context *CommandContext) {},
	}
}
