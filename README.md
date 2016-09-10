# BT Home Hub CLI

A CLI that can interact with BT Home Hub routers. So far it's only been proven against
hub version `Home Hub 60 Type A` running firmware `SG4B100021EC`.

[![asciicast](https://asciinema.org/a/4u35xe98mgj1lc7olrl1v0gie.png)](https://asciinema.org/a/4u35xe98mgj1lc7olrl1v0gie)

Useful for performing quick reboots or for pulling statistics from the Home Hub.

## Usage

#### Using the CLI

Download one of the [releases](https://github.com/jamesnetherton/homehub-cli/releases) for your operating system.

For *nix and OS X:

    ./homehub-cli

For Windows:

    homehub-cli.exe

You'll need to authenticate against the Home Hub in order to do anything useful. Use the `Login` command to do this. Pressing the `TAB` key shows all of the available commands.

#### Running indivdual commands

You can run indivdual commands outside of the CLI shell by specifying the desired function to execute, together with the Home Hub authentication details.

    ./homehub-cli Reboot --huburl=http://192.168.1.254 --username=admin --password=bar

The `huburl` and `username` arguments are defaulted to the standard Home Hub IP address and admin user name. So you can omit these arguments if you want to and just specify the `password`.

See `./homehub-cli --help` for all options.

## Building

    git clone git@github.com:jamesnetherton/homehub-cli.git $GOPATH/src/github.com/jamesnetherton/homehub-cli
    make build

Generated binaries are output to the `build` directory.
