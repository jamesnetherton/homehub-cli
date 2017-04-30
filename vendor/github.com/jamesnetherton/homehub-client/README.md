# Home Hub Client

A golang client that can interact with BT Home Hub routers. Refer to the [compatibility matrix](matrix.md)
to see the firmware versions supported by each release. The master branch is currently proven against firmware version `SG4B10002244`.

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
