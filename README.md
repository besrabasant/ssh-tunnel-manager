## SSH Tunnel Manager (sshtm)

SSH Tunnel Manager (**sshtm**) allows you to manage SSH tunnel configurations for easy setup and use of port forwarding via SSH.

## Current Version

**SSH Tunnel Manager (sshtm) v1.0.0**

Ensure you are using the latest version to benefit from the newest features and improvements.


## Features

- **Save Configuration**: Store SSH tunnel configurations including descriptions for later use.
- **Start Tunnel**: Initiate an SSH tunnel using a saved configuration, similar to the SSH command `ssh -L [LOCAL_IP:]LOCAL_PORT:DESTINATION:DESTINATION_PORT [USER@]SSH_SERVER`.
- **List Configurations**: Display all saved SSH tunnel configurations.
- **Remove Configuration**: Delete a saved SSH tunnel configuration.

## Requirements

To build and run sshtm, you must have the following installed on your system:

- **Go (Golang)**: Ensure you have Go installed (version 1.15 or later is recommended). You can check your current Go version using `go version`.
- **Git**: Needed for cloning the repository and potentially fetching dependencies.
- **Air (optional for development)**: For live reloading during development. Install Air by following the instructions at [Air GitHub Repository](https://github.com/cosmtrek/air).



## Installation
Clone the repository and navigate to the project directory:

```sh
git clone https://github.com/besrabasant/ssh-tunnel-manager.git
cd ssh-tunnel-manager
```

To build and install **sshtm**, run:

```sh
./install.sh
```
This script will compile the project, install executable files **sshtm**, and set up necessary environment variables.

## Uninstallation
```sh
~/.local/share/sshtm/scripts/uninstall.sh
```

## Basic Usage

```
sshtm [command]
```

## Commands

Below is the list of available commands in **sshtm**, with descriptions for each:

| Command        | Aliases             | Description                                                                                                                                                         |
| -------------- | ------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **list**       | **ls**, **l**       | List all SSH tunnel configurations, optionally filtering by a search pattern. (You can use a pattern to only list the configurations that fuzzy match that pattern) |  |
| **add**        | **a**               | Add a new SSH tunnel configuration.                                                                                                                                 |
| **edit**       | **e**               | Edit an existing SSH tunnel configuration.                                                                                                                          |
| **delete**     | **del**, **d**      | Delete an existing SSH tunnel configuration.                                                                                                                        |
| **tunnel**     | **t**               | Start an SSH tunnel using a saved configuration, optionally specifying a local port.                                                                                |
| **active**     |                     | List all active SSH tunnels.                                                                                                                                        |
| **kill**       | **k**,**terminate** | Terminate an active SSH tunnel.                                                                                                                                     |
| **help**       |                     | Display help about any command.                                                                                                                                     |
| **completion** |                     | Generate the autocompletion script for the specified shell.                                                                                                         |
| **version**    |                     | Print the version number of SSH Tunnel Manager                                                                                                                      |

## Development

1. **Start daemon** :
   To start the sshtm background service for development, run:
   ```sh
    make gen_proto
    air
   ```

   This command assumes you have [Air](https://github.com/cosmtrek/air) installed for live reloading during development. It should be executed from the project root directory.
   
2. **Use the CLI** :
   To interact with sshtm via the command line interface during development, navigate to the `client` directory:
   ```sh
   cd client
   go run main.go [command] [arguments]
   ```
   Here, **[command]** is one of the commands available and **[arguments]** represents additional arguments that can be passed to the command.


## Versioning Policy

This project uses [Semantic Versioning (SemVer)](https://semver.org/) for version numbers.
