package tasks

import (
	"context"
	"fmt"
	"strings"

	"github.com/besrabasant/ssh-tunnel-manager/rpc"
)

func KillTunnelTask(ctx context.Context, req *rpc.KillTunnelRequest, manager *TunnelManager) (*rpc.KillTunnelResponse, error) {
	var output strings.Builder

	output.WriteString("\n")

	if connInfo, exists := manager.connections[int(req.LocalPort)]; exists {
		output.WriteString(fmt.Sprint("Closing existing connection on port ", int(req.LocalPort), "\n"))

		// If there's an existing connection on the same port, close it
		connInfo.Cancel() // Cancel the context of the existing connection

		output.WriteString(fmt.Sprintf("\nTunneling stopped %q <==> %q through %q\n", connInfo.LocalAddr, connInfo.RemoteAddr, connInfo.Config.Server))

	} else {
		output.WriteString(fmt.Sprint("Did not find any connection on port ", int(req.LocalPort), "\n"))
	}

	return &rpc.KillTunnelResponse{Result: output.String()}, nil
}
