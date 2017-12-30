# Home Hub CLI

A CLI that can interact with BT Home Hub routers.

[![asciicast](https://asciinema.org/a/4u35xe98mgj1lc7olrl1v0gie.png)](https://asciinema.org/a/4u35xe98mgj1lc7olrl1v0gie)

Useful for performing quick reboots or for pulling statistics from the Home Hub.

## Usage

#### Using the CLI

Download one of the [releases](https://github.com/jamesnetherton/homehub-cli/releases) for your operating system. Please read the release notes to ensure
it's compatible with your Home Hub firmware version. Or refer to the [compatibility matrix](matrix.md).

For *nix and OS X:

```
./homehub-cli
```

For Windows:

```
homehub-cli.exe
```

Pressing the `TAB` key shows all of the available commands. You can get help for specific commands by doing:

`CommandName --help`

#### Running indivdual commands

You can run indivdual commands outside of the CLI shell by specifying the desired function to execute, together with the Home Hub authentication details.

```
./homehub-cli Reboot --huburl=http://192.168.1.254 --username=admin --password=secret
```

If the command requires arguments, you can pass them along like this.

```
./homehub-cli WiFiFrequency24GhzChannelSet 6 --huburl=http://192.168.1.254 --username=admin --password=secret
```

The `huburl` and `username` arguments are defaulted to the standard Home Hub IP address and admin user name. So you can omit these arguments if you want to and just specify the `password`.

See `./homehub-cli --help` for all options.

## Troubleshooting

### Compatible firmware versions

First check that the CLI is [compatible](matrix.md) with your router hardware. You can run the `About` command to see which firmware version the homehub-cli was tested and verified against.

### Enable debug logging

You can enable debugging within the CLI via the `EnableDebug` command.

To enable debugging:

`EnableDebug true`

To disable debugging:

`EnableDebug false`

With debugging enabled, you'll see HTTP requests and responses being output whenever commands are executed.

## Building

    git clone git@github.com:jamesnetherton/homehub-cli.git $GOPATH/src/github.com/jamesnetherton/homehub-cli
    make build

Generated binaries are output to the `build` directory.
