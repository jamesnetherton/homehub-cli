package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/jamesnetherton/homehub-client"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/cobra/cmd"
)

var (
	hubURL   string
	userName string
	password string
	hub      *homehub.Hub
)

func main() {
	cmd.RootCmd.PersistentFlags().StringVarP(&hubURL, "huburl", "r", "http://192.168.1.254", "URL of the home hub router")
	cmd.RootCmd.PersistentFlags().StringVarP(&userName, "username", "u", "admin", "The hub router user name")
	cmd.RootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "The home hub router password")

	cmdHandler := func(cmd *cobra.Command, args []string) {
		if !stringIsEmpty(hubURL) && !stringIsEmpty(userName) && !stringIsEmpty(password) {
			hub = homehub.New(hubURL, userName, password)
			success, _ := hub.Login()

			if !success {
				fmt.Println("Login failed")
				return
			}

			invokeMethod(nil, cmd.Name(), args)
		} else {
			cmd.Usage()
		}
	}

	for _, funcName := range getFuncNames() {
		cmd.RootCmd.AddCommand(&cobra.Command{
			Use: funcName,
			Run: cmdHandler,
		})
	}

	helpFunc := func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage:\n  homehub-cli [command] --huburl=<home hub url> --username=<home hub username> --password=<home hub password>")
		fmt.Println("\nCommands:\n ", strings.Join(getFuncNames(), "\n  "))
	}

	usageFunc := func(cmd *cobra.Command) error {
		fmt.Printf("Usage:\n  homehub-cli %s --huburl=<home hub url> --username=<home hub username> --password=<home hub password>\n", cmd.Name())
		return nil
	}

	cmd.RootCmd.SetHelpFunc(helpFunc)
	cmd.RootCmd.SetUsageFunc(usageFunc)
	cmd.RootCmd.SilenceErrors = true

	if len(os.Args[1:]) == 0 {
		banner()

		l, err := createReadline()
		if err != nil {
			panic(err)
		}
		defer l.Close()

		for {
			line, err := l.Readline()

			if err == readline.ErrInterrupt {
				if len(line) == 0 {
					break
				} else {
					continue
				}
			} else if err == io.EOF {
				break
			}

			line = strings.TrimSpace(line)
			if !stringIsEmpty(line) {
				if !strings.HasPrefix(line, "Login") {
					cmd := strings.Split(line, " ")
					methodName := cmd[0]
					args := append(cmd[:0], cmd[1:]...)
					invokeMethod(l, methodName, args)
				} else {
					hub = doLogin(l)
				}
			}
		}
	} else {
		err := cmd.RootCmd.Execute()
		if err != nil {
			helpFunc(nil, nil)
		}
	}
}

func createReadline() (l *readline.Instance, err error) {
	return readline.NewEx(&readline.Config{
		Prompt:          getUserPrompt(),
		AutoComplete:    readline.NewPrefixCompleter(readline.PcItemDynamic(listFuncNames())),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
}

func getUserPrompt() string {
	var user string

	if !stringIsEmpty(os.Getenv("USER")) {
		user = os.Getenv("USER")
	} else if !stringIsEmpty(os.Getenv("USERNAME")) {
		user = os.Getenv("USERNAME")
	} else {
		user = "unknown"
	}

	return fmt.Sprintf("%s@homehub: ", user)
}

func listFuncNames() func(string) []string {
	return func(s string) []string {
		return getFuncNames()
	}
}

func getFuncNames() []string {
	v := reflect.TypeOf(&homehub.Hub{})
	funcNames := make([]string, v.NumMethod())
	for i := 0; i < v.NumMethod(); i++ {
		funcNames[i] = v.Method(i).Name
	}
	return funcNames
}

func invokeMethod(l *readline.Instance, methodName string, args []string) {
	h := reflect.TypeOf(hub)
	m, found := h.MethodByName(methodName)
	if found {
		t := m.Func.Type()
		inputs := make([]reflect.Value, t.NumIn())

		if t.NumIn() > 0 {
			argOffset := 0
			if t.In(0).Kind().String() == "ptr" {
				argOffset = 1
			}

			if len(args) != t.NumIn()-argOffset {
				fmt.Printf("Wrong number of arguments for '%s'. Expected %d but got %d.\n", methodName, t.NumIn()-1, len(args))
				return
			}

			for i := argOffset; i < t.NumIn(); i++ {
				switch t.In(i).Kind().String() {
				case "bool":
					val, err := strconv.ParseBool(args[i-argOffset])
					if err == nil {
						inputs[i] = reflect.ValueOf(val)
					} else {
						fmt.Printf("Expected '%s' boolean function argument %d to be either true or false . But got '%s'.\n", methodName, i, args[i-argOffset])
						return
					}
					break
				case "int":
					val, err := strconv.Atoi(args[i-argOffset])
					if err == nil {
						inputs[i] = reflect.ValueOf(val)
					} else {
						fmt.Printf("Expected '%s' int function argument %d to be numeric. But got '%s'.\n", methodName, i, args[i-argOffset])
						return
					}
					break
				default:
					inputs[i] = reflect.ValueOf(args[i-argOffset])
				}
			}
		}

		if hub == nil {
			fmt.Printf("\nYou are not logged in. Please login...\n\n")
			hub = doLogin(l)
		}

		if hub != nil {
			inputs[0] = reflect.ValueOf(hub)
			resp := m.Func.Call(inputs)

			if t.NumOut() > 0 {
				err := ""
				for i := 0; i < t.NumOut(); i++ {
					if resp[i].Type().String() == "error" && !resp[i].IsNil() {
						err = fmt.Sprintf("%s", resp[i].Interface())
						break
					}
				}
				if len(err) > 0 {
					if err == "Invalid user session" {
						fmt.Println("Login session has expired")
						hub = nil
						invokeMethod(l, methodName, args)
					} else {
						fmt.Printf("Error executing command '%s': %s\n", methodName, err)
					}
				} else if resp[0].Type().String() != "error" {
					fmt.Println(resp[0].String())
				}
			}
		}
	} else {
		fmt.Println("homehub: Command not found:", strconv.Quote(methodName))
	}
}

func stringIsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func doLogin(l *readline.Instance) *homehub.Hub {
	var URL string
	var username string
	var password []byte

	fmt.Print("Home hub URL: ")
	fmt.Scan(&URL)

	fmt.Print("Home hub user name: ")
	fmt.Scan(&username)

	password, _ = l.ReadPassword("Home hub password: ")

	hub := homehub.New(URL, username, string(password))
	success, _ := hub.Login()

	if !success {
		fmt.Println("Login failed")
		return nil
	}

	fmt.Println("Logged in as", username)

	return hub
}

func banner() {
	fmt.Println(" _   _                           _   _         _")
	fmt.Println("| | | |                         | | | |       | |")
	fmt.Println("| |_| |  ___   _ __ ___    ___  | |_| | _   _ | |__")
	fmt.Println("|  _  | / _ \\ | '_ ` _ \\  / _ \\ |  _  || | | || '_ \\")
	fmt.Println("| | | || (_) || | | | | ||  __/ | | | || |_| || |_) |")
	fmt.Println("\\_| |_/ \\___/ |_| |_| |_| \\___| \\_| |_/ \\__,_||_.__/")
	fmt.Printf("\n=====================================================\n\n")
}
