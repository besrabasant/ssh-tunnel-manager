package tasks

import (
	"strings"

	"github.com/besrabasant/ssh-tunnel-manager/rpc"
)

func mutationResponse(result string, status rpc.ResponseStatus, message string) *rpc.AddOrUpdateConfigurationResponse {
	return &rpc.AddOrUpdateConfigurationResponse{
		Result:  result,
		Status:  status,
		Message: message,
	}
}

func operationParts(result string) (rpc.ResponseStatus, []string) {
	trimmed := strings.TrimSpace(result)
	if trimmed == "" {
		return rpc.ResponseStatus_Success, nil
	}

	events := strings.Split(trimmed, "\n")
	status := rpc.ResponseStatus_Success
	for _, event := range events {
		lower := strings.ToLower(event)
		if strings.Contains(lower, "failed") || strings.Contains(lower, "error") {
			status = rpc.ResponseStatus_Error
			break
		}
	}
	return status, events
}
