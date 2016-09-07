# BT Home Hub Client

A golang client that can interact with BT Home Hub routers. So far it's only been proven against
hub version `Home Hub 60 Type A` running firmware `SG4B100021EC`.

At present, only a small set of the available [APIs](xpath.go) have been implemented.

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
