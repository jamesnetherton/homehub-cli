# BT Home Hub CLI

A CLI that can interact with BT Home Hub routers. So far it's only been proven against
hub version `Home Hub 60 Type A` running firmware `SG4B100021EC`.

[![asciicast](https://asciinema.org/a/4u35xe98mgj1lc7olrl1v0gie.png)](https://asciinema.org/a/4u35xe98mgj1lc7olrl1v0gie)

Useful for performing quick reboots or for pulling statistics from the Home Hub.

## Usage

Download one of the [releases](releases) for your operating system.

For *nix and OS X:

    ./homehub-cli

For Windows:

    homehub-cli.exe

You'll need to authenticate against the Home Hub in order to do anything useful. Use the `Login` command to do this. Pressing the `TAB` key shows all of the available commands.

## Building

    git clone git@github.com:jamesnetherton/homehub-cli.git $GOPATH/src/github.com/jamesnetherton/homehub-cli
    make build

Generated binaries are output to the `build` directory.
