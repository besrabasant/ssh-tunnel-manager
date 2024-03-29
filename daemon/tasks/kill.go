package tasks

import (
	"context"
	"fmt"
	"strings"

	"github.com/besrabasant/ssh-tunnel-manager/rpc"
)

func KillTunnelTask(ctx context.Context, req *rpc.KillTunnelRequest) (*rpc.KillTunnelResponse, error) {
	var output strings.Builder

	output.WriteString("\n")

	ConnMutex.Lock()
	if connInfo, exists := Connections[int(req.LocalPort)]; exists {
		output.WriteString(fmt.Sprint("Closing existing connection on port ", int(req.LocalPort), "\n"))

		// If there's an existing connection on the same port, close it
		connInfo.Cancel() // Cancel the context of the existing connection
		if connInfo.Listen != nil {
			connInfo.StopListening()
			connInfo.KillAllConns()
		}
		output.WriteString(fmt.Sprint("Tunnel stopped", int(req.LocalPort), "\n"))
	} else {
		output.WriteString(fmt.Sprint("Did not find any connection on port ", int(req.LocalPort), "\n"))
	}
	ConnMutex.Unlock()

	return &rpc.KillTunnelResponse{Result: output.String()}, nil
}
