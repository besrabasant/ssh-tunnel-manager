## SSH Tunnel Manager (sshtm)

SSH Tunnel Manager (**sshtm**) allows you to manage SSH tunnel configurations for easy setup and use of port forwarding via SSH.

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

| Command    | Description                                                                                                                          |
| ---------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| **list**       | List all configurations. Supports fuzzy matching to filter configurations by a pattern.                                              |
| **add**        | Add a new SSH tunnel configuration.                                                                                                  |
| **edit**       | Edit an existing SSH tunnel configuration.                                                                                           |
| **delete**     | Delete an existing SSH tunnel configuration.                                                                                         |
| **help**       | Display help about any command.                                                                                                      |
| **tunnel**     | Start an SSH tunnel using a saved configuration. Forwards connections to the specified local port or a random port if not specified. |
| **active**     | List all active SSH tunnels.                                                                                                         |
| **kill**       | Terminate an active SSH tunnel.                                                                                                      |
| **completion** | Generate the autocompletion script for the specified shell.                                                                          |


## Development

1. **Start daemon** :
   To start the sshtm background service for development, run:
    ```sh
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