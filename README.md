## SSH tunnel manager

Save SSH tunnel configurations and start a tunnel with port forwarding using one of the saved configurations.

## Features

- Save an SSH tunnel configuration with a description
- Start an SSH tunnel with port forwarding using a configuration. Same as `ssh -L [LOCAL_IP:]LOCAL_PORT:DESTINATION:DESTINATION_PORT [USER@]SSH_SERVER`
- List all configurations
- Remove a configuration

## Requirements


## Build

```sh
make

./install.sh
```
## Usage

```
sshtm [command]
```


## Development Usage

1. **Start daemon** :
   Run from project root
    ```sh
    air
    ```
2. **Use Cli** :
   Run inside **client** directory
   ```sh
   go run main.go [command] [--opts]
   ```