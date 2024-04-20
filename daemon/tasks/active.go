package tasks

import (
	"context"
	"fmt"
	"strings"

	"github.com/besrabasant/ssh-tunnel-manager/pkg/tunnelmanager"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/besrabasant/ssh-tunnel-manager/utils"
)

func ListActiveTunnelsTask(ctx context.Context, req *rpc.ListActiveTunnelsRequest, manager *tunnelmanager.TunnelManager) (*rpc.ListActiveTunnelsResponse, error) {
	var output strings.Builder

	output.WriteString("\n")

	if len(manager.Connections) > 0 {
		for port, conn := range manager.Connections {
			// config is prented without a new line at its end.
			writeTunnelsToOutput(&output, port, &conn)
			output.WriteString("\n")
		}
	} else {
		output.WriteString("\nNo active connections found.\n")
	}

	return &rpc.ListActiveTunnelsResponse{Result: output.String()}, nil
}

func writeTunnelsToOutput(out *strings.Builder, port int, conn *tunnelmanager.ConnectionInfo) {
	template := `:%s
 - Connection:     		%s
 - Remote Address: 		%s
 - Local Address:		%s
 - SSH server:			%s
`
	portStr := utils.Bold(fmt.Sprint(port))
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
