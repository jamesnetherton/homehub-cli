package cmd

import (
	"fmt"
	"net/url"
	"strings"

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
				hubURL, _ := context.ReadLine("Home Hub URL: ")
				if len(strings.TrimSpace(hubURL)) == 0 {
					fmt.Println("Hub URL must not be empty")
					return
				}

				_, err := url.ParseRequestURI(hubURL)
				if err != nil {
					fmt.Println("Hub URL must be a valid URL")
					return
				}

				userName, _ := context.ReadLine("Home Hub user: ")
				if len(strings.TrimSpace(userName)) == 0 {
					fmt.Println("Hub user must not be empty")
					return
				}

				password, _ := context.ReadPassword("Home Hub password: ")
				if len(strings.TrimSpace(password)) == 0 {
					fmt.Println("Hub password must not be empty")
					return
				}

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
