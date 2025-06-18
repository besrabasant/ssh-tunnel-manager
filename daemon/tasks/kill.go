package tasks

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/besrabasant/ssh-tunnel-manager/config"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/tunnelmanager"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/besrabasant/ssh-tunnel-manager/utils"
)

func KillTunnelTask(ctx context.Context, req *rpc.KillTunnelRequest, manager *tunnelmanager.TunnelManager) (*rpc.KillTunnelResponse, error) {
	var output strings.Builder
	var connInfo *tunnelmanager.ConnectionInfo
	var connPort int

	output.WriteString("\n")

	if req.ConfigName != "" {
		openConns := manager.Connections.Filter(func(c *tunnelmanager.ConnectionInfo) bool {
			return req.ConfigName == c.Config.Name
		})

		if len(openConns) > 0 {
			port, conn, found := utils.GetFirstItemFromMap(openConns)
			if found {
				connInfo = conn
				connPort = port
			}
		}
	} else {
		connPort = int(req.LocalPort)
		connInfo = manager.Connections[connPort]
	}

	if connInfo != nil {
		manager.Mutex.Lock()
		defer manager.Mutex.Unlock()

		output.WriteString(fmt.Sprint("Closing existing connection on port ", connPort, " for ", connInfo.Config.Name, "\n"))

		// If there's an existing connection on the same port, close it
		connInfo.ClearConnection() // Cancel the context of the existing connection

		output.WriteString(fmt.Sprintf("\nTunneling stopped %q <==> %q through %q\n", connInfo.LocalAddr, connInfo.RemoteAddr, connInfo.Config.Server))

	} else {
		if req.ConfigName != "" {
			output.WriteString(fmt.Sprint("Did not find any connection for configuration ", req.ConfigName, "\n"))
		} else {
			output.WriteString(fmt.Sprint("Did not find any connection on port ", connPort, "\n"))
		}
	}

	dirpath := config.DefaultConfigDir
	if value := os.Getenv(config.ConfigDirFlagName); value != "" {
		dirpath = value
	}

	configdir, err := utils.ResolveDir(dirpath)
	if err == nil {
		manager.SaveActiveTunnels(filepath.Join(configdir, config.ActiveTunnelsFile))
	}

	return &rpc.KillTunnelResponse{Result: output.String()}, nil
}
