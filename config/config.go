package config

var AppVersion = "1.0.5"

const DefaultSSHPort = "22"

const Host = "localhost"

const Port = "50051"

const Address = Host + ":" + Port

const ConfigDirFlagName = "config-dir"

const DefaultConfigDir = "~/.ssh-tunnel-manager"

const ActiveTunnelsFile = "active_tunnels.json"
