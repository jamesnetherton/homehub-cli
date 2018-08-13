# Home Hub Client

[![CircleCI](https://img.shields.io/circleci/project/jamesnetherton/homehub-client/master.svg)](https://circleci.com/gh/jamesnetherton/homehub-client/tree/master)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=600)](https://opensource.org/licenses/MIT)

A golang client that can interact with BT Home Hub routers. Refer to the [compatibility matrix](matrix.md)
to see the firmware versions supported by each release. The master branch is currently proven against firmware versions `SG4B1000B540` and `SG4B1A006100`.

At present, only a small set of the available [APIs](xpath.go) have been implemented.

If you're looking for a command line implementation of this library, check out my [Home Hub CLI](https://github.com/jamesnetherton/homehub-cli).

## Usage

```golang
package main

import (
	"fmt"
	"github.com/jamesnetherton/homehub-client"
)

func main() {
	hub := homehub.New("http://192.168.1.254", "admin", "p4ssw0rd")
	hub.Login()

	version, err := hub.Version()
	if err == nil {
		fmt.Printf("Home Hub Version = %s", version)
	}
}
```
