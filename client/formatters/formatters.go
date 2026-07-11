// Package formatters contains the terminal presentation layer for sshtm.
package formatters

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/charmbracelet/lipgloss"
)

var ansiEscape = regexp.MustCompile("\\x1b\\[[0-9;]*m")

// Formatter renders a typed value for a human-readable CLI command.
type Formatter[T any] interface {
	Format(T) string
}

type styles struct {
	title   lipgloss.Style
	rule    lipgloss.Style
	heading lipgloss.Style
	label   lipgloss.Style
	value   lipgloss.Style
	success lipgloss.Style
	error   lipgloss.Style
	muted   lipgloss.Style
}

func newStyles(w io.Writer) styles {
	renderer := lipgloss.NewRenderer(w)
	return styles{
		title:   renderer.NewStyle().Bold(true).Foreground(lipgloss.Color("205")),
		rule:    renderer.NewStyle().Foreground(lipgloss.Color("240")),
		heading: renderer.NewStyle().Bold(true).Foreground(lipgloss.Color("81")),
		label:   renderer.NewStyle().Foreground(lipgloss.Color("245")),
		value:   renderer.NewStyle().Foreground(lipgloss.Color("252")),
		success: renderer.NewStyle().Foreground(lipgloss.Color("42")),
		error:   renderer.NewStyle().Foreground(lipgloss.Color("203")),
		muted:   renderer.NewStyle().Foreground(lipgloss.Color("245")),
	}
}

func writeLine(builder *strings.Builder, line string) {
	builder.WriteString(line)
	builder.WriteByte('\n')
}

// ConfigurationListFormatter renders saved tunnel configurations.
type ConfigurationListFormatter struct {
	styles styles
}

func NewConfigurationListFormatter(w io.Writer) Formatter[*rpc.ListConfigurationsResponse] {
	return ConfigurationListFormatter{styles: newStyles(w)}
}

func (f ConfigurationListFormatter) Format(response *rpc.ListConfigurationsResponse) string {
	if response == nil {
		return ""
	}
	if len(response.GetConfigs()) == 0 {
		return legacyText(response.GetResult())
	}

	var output strings.Builder
	for i, cfg := range response.GetConfigs() {
		if i > 0 {
			output.WriteByte('\n')
		}
		writeLine(&output, f.styles.heading.Render(":"+formatPort(cfg.GetLocalPort())))
		writeConnectionProperty(&output, f.styles, "Connection:", cfg.GetName()+descriptionSuffix(cfg.GetDescription(), f.styles.muted))
		writeConnectionProperty(&output, f.styles, "Remote Address:", fmt.Sprintf("%s:%d", cfg.GetRemoteHost(), cfg.GetRemotePort()))
		writeConnectionProperty(&output, f.styles, "Local Address:", localAddress(cfg.GetLocalPort()))
		writeConnectionProperty(&output, f.styles, "SSH server:", cfg.GetServer())
		writeConnectionProperty(&output, f.styles, "User:", cfg.GetUser())
		writeConnectionProperty(&output, f.styles, "Private key:", cfg.GetKeyFile())
	}
	return output.String()
}

func descriptionSuffix(description string, muted lipgloss.Style) string {
	if strings.TrimSpace(description) == "" {
		return ""
	}
	return " " + muted.Render("("+description+")")
}

func formatPort(port int32) string {
	if port == 0 {
		return "auto"
	}
	return fmt.Sprintf("%d", port)
}

func localAddress(port int32) string {
	if port == 0 {
		return "auto"
	}
	return fmt.Sprintf("127.0.0.1:%d", port)
}

// ActiveTunnelsFormatter renders currently active tunnels.
type ActiveTunnelsFormatter struct {
	styles styles
}

func NewActiveTunnelsFormatter(w io.Writer) Formatter[*rpc.ListActiveTunnelsResponse] {
	return ActiveTunnelsFormatter{styles: newStyles(w)}
}

func (f ActiveTunnelsFormatter) Format(response *rpc.ListActiveTunnelsResponse) string {
	if response == nil {
		return ""
	}
	if len(response.GetTunnels()) == 0 {
		return legacyText(response.GetResult())
	}

	var output strings.Builder
	for i, tunnel := range response.GetTunnels() {
		if i > 0 {
			output.WriteByte('\n')
		}
		writeLine(&output, f.styles.heading.Render(fmt.Sprintf(":%d", tunnel.GetLocalPort())))
		writeConnectionProperty(&output, f.styles, "Connection:", tunnel.GetName())
		writeConnectionProperty(&output, f.styles, "Remote Address:", tunnel.GetRemoteAddr())
		localAddress := tunnel.GetLocalAddr()
		if localAddress == "" {
			localAddress = localAddressForPort(tunnel.GetLocalPort())
		}
		writeConnectionProperty(&output, f.styles, "Local Address:", localAddress)
		writeConnectionProperty(&output, f.styles, "SSH server:", tunnel.GetServer())
	}
	return output.String()
}

func localAddressForPort(port int32) string {
	return fmt.Sprintf("127.0.0.1:%d", port)
}

func writeConnectionProperty(builder *strings.Builder, s styles, label, value string) {
	writeLine(builder, s.label.Render(fmt.Sprintf("- %-29s", label))+s.value.Render(value))
}

// MutationFormatter renders add, update, delete, and kill responses.
type MutationFormatter struct {
	styles styles
}

func NewMutationFormatter(w io.Writer) MutationFormatter {
	return MutationFormatter{styles: newStyles(w)}
}

func (f MutationFormatter) Format(result MutationResult) string {
	if result.Structured {
		if result.Status == rpc.ResponseStatus_Error {
			return f.styles.error.Render(result.Message) + "\n"
		}
		return f.styles.success.Render(result.Message) + "\n"
	}
	return legacyText(result.Legacy)
}

// MutationResult is the common client-side representation of a mutation.
type MutationResult struct {
	Status     rpc.ResponseStatus
	Message    string
	Legacy     string
	Structured bool
}

func MutationFromAddOrUpdate(response *rpc.AddOrUpdateConfigurationResponse) MutationResult {
	if response == nil {
		return MutationResult{}
	}
	return MutationResult{
		Status:     response.GetStatus(),
		Message:    response.GetMessage(),
		Legacy:     response.GetResult(),
		Structured: response.GetMessage() != "" || response.GetData() != nil,
	}
}

func MutationFromDelete(response *rpc.DeleteConfigurationResponse) MutationResult {
	if response == nil {
		return MutationResult{}
	}
	return MutationResult{
		Status:     response.GetStatus(),
		Message:    response.GetMessage(),
		Legacy:     response.GetResult(),
		Structured: response.GetMessage() != "",
	}
}

func MutationFromKill(response *rpc.KillTunnelResponse) MutationResult {
	if response == nil {
		return MutationResult{}
	}
	return MutationResult{
		Status:     response.GetStatus(),
		Message:    response.GetMessage(),
		Legacy:     response.GetResult(),
		Structured: response.GetMessage() != "",
	}
}

// OperationFormatter renders tunnel start progress and operation results.
type OperationFormatter struct {
	styles styles
}

func NewOperationFormatter(w io.Writer) OperationFormatter {
	return OperationFormatter{styles: newStyles(w)}
}

func (f OperationFormatter) Format(result OperationResult) string {
	if !result.Structured {
		return legacyText(result.Legacy)
	}

	var output strings.Builder
	for _, event := range result.Events {
		if strings.Contains(strings.ToLower(event), "failed") || strings.Contains(strings.ToLower(event), "error") {
			writeLine(&output, f.styles.error.Render(event))
		} else {
			writeLine(&output, f.styles.value.Render(event))
		}
	}
	if result.Message != "" {
		style := f.styles.success
		if result.Status == rpc.ResponseStatus_Error {
			style = f.styles.error
		}
		writeLine(&output, style.Render(result.Message))
	}
	return output.String()
}

func legacyText(value string) string {
	return ansiEscape.ReplaceAllString(value, "")
}

// OperationResult is the common client-side representation of a tunnel operation.
type OperationResult struct {
	Status     rpc.ResponseStatus
	Message    string
	Events     []string
	Legacy     string
	Structured bool
}

func OperationFromStart(response *rpc.StartTunnelResponse) OperationResult {
	if response == nil {
		return OperationResult{}
	}
	return OperationResult{
		Status:     response.GetStatus(),
		Message:    response.GetMessage(),
		Events:     response.GetEvents(),
		Legacy:     response.GetResult(),
		Structured: len(response.GetEvents()) > 0 || response.GetMessage() != "",
	}
}

// FetchFormatter renders fetch errors used by the edit command.
type FetchFormatter struct {
	styles styles
}

func NewFetchFormatter(w io.Writer) FetchFormatter {
	return FetchFormatter{styles: newStyles(w)}
}

func (f FetchFormatter) Format(response *rpc.FetchConfigurationResponse) string {
	if response == nil {
		return ""
	}
	if response.GetStatus() == rpc.ResponseStatus_Error {
		return f.styles.error.Render(response.GetMessage()) + "\n"
	}
	return response.GetMessage()
}
