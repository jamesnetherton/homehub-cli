package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/jamesnetherton/homehub-cli/functions"
	"github.com/jamesnetherton/homehub-client"
)

func main() {
	var hub *homehub.Hub
	var user string

	if len(strings.TrimSpace(os.Getenv("USER"))) > 0 {
		user = os.Getenv("USER")
	} else if len(strings.TrimSpace(os.Getenv("USERNAME"))) > 0 {
		user = os.Getenv("USERNAME")
	} else {
		user = "unknown"
	}

	banner()
	createCompleter()

	completer := createCompleter()
	l, err := readline.NewEx(&readline.Config{
		Prompt:          user + "@homehub: ",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
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
		if len(line) > 0 {
			if line != "Login" {
				invokeMethod(hub, line)
			} else {
				hub = doLogin(l)
			}
		}
	}
}

func createCompleter() *readline.PrefixCompleter {
	return readline.NewPrefixCompleter(readline.PcItemDynamic(listFuncNames()))
}

func listFuncNames() func(string) []string {
	return func(s string) []string {
		return functions.FuncNames
	}
}

func invokeMethod(hub *homehub.Hub, methodName string) {
	m := reflect.ValueOf(hub).MethodByName(methodName)
	if m.IsValid() {
		if hub != nil {
			resp := m.Call(nil)
			result := resp[0].String()
			fmt.Println(result)
		} else {
			fmt.Println("homehub: You are not logged in. Use the Login command")
		}
	} else {
		fmt.Println("homehub: Command not found:", strconv.Quote(methodName))
	}
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
