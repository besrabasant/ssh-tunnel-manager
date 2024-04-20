package configmanager

import (
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
)

func ConvertConfigToRpcTunnelConfig(cfg *Entry) *rpc.TunnelConfig {
	return &rpc.TunnelConfig{
		Name:        cfg.Name,
		Description: cfg.Description,
		Server:      cfg.Server,
		User:        cfg.User,
		KeyFile:     cfg.KeyFile,
		RemoteHost:  cfg.RemoteHost,
		RemotePort:  int32(cfg.RemotePort),
		LocalPort:  int32(cfg.LocalPort),
	}
}

func ConvertRpcTunnelConfigToConfig(cfg *rpc.TunnelConfig) *Entry {
	return &Entry{
		Name:        cfg.Name,
		Description: cfg.Description,
		Server:      cfg.Server,
		User:        cfg.User,
		KeyFile:     cfg.KeyFile,
		RemoteHost:  cfg.RemoteHost,
		RemotePort:  int(cfg.RemotePort),
		LocalPort:  int(cfg.LocalPort),
	}
}