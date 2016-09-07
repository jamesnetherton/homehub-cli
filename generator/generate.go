package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	response, err := http.Get("https://raw.githubusercontent.com/jamesnetherton/homehub-client/master/hub.go")
	if err == nil {
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		body := string(bodyBytes[:])

		output := "package functions\n\nvar FuncNames = []string {\n"
		var names []string

		regex := regexp.MustCompile("func \\(h \\*Hub\\) (.*?[^\\(]*)")

		for _, line := range strings.Split(body, "\n") {
			if strings.HasPrefix(line, "func") {
				matches := [][]string{regex.FindStringSubmatch(line)}
				if len(matches[0]) > 0 {
					names = append(names, strconv.Quote(matches[0][1]))
				}
			}
		}
		output += strings.Join(names, ",\n")
		output += "}"

		ioutil.WriteFile("../functions/funcs.go", []byte(output), 0644)

		cmd := exec.Command("gofmt", "-w", "func_names.go")
		fmt.Println(cmd.Run())
	}
}
