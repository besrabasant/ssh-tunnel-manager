package tasks

import (
	"context"
	"fmt"
	"strings"

	"github.com/besrabasant/ssh-tunnel-manager/rpc"
)

func ListActiveTunnelsTask(ctx context.Context, req *rpc.ListActiveTunnelsRequest, manager *TunnelManager) (*rpc.ListActiveTunnelsResponse, error) {
	var output strings.Builder

	output.WriteString("\n")

	if len(manager.connections) > 0 {
		for port, conn := range manager.connections {
			// config is prented without a new line at its end.
			writeTunnesToOutput(&output, port, &conn)
			output.WriteString("\n")
		}
	} else {
		output.WriteString("\nNo active connections found.\n")
	}

	return &rpc.ListActiveTunnelsResponse{Result: output.String()}, nil
}

func writeTunnesToOutput(out *strings.Builder, port int, conn *ConnectionInfo) {
	template := `:%s
 - Connection:     		%s
 - Remote Address: 		%s
 - Local Address:		%s
 - SSH server:			%s
`
	portStr := bold(fmt.Sprint(port))
	out.Write([]byte(
		fmt.Sprintf(
			template,
			portStr,
			conn.Config.Name,
			conn.RemoteAddr,
			conn.LocalAddr,
			conn.Config.Server,
		)))
}
