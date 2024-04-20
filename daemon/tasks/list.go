package tasks

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/besrabasant/ssh-tunnel-manager/config"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
	pb "github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/besrabasant/ssh-tunnel-manager/utils"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

func getConfigs() ([]configmanager.Entry, error) {
	dirpath := config.DefaultConfigDir
	if value := os.Getenv(config.ConfigDirFlagName); value != "" {
		dirpath = value
	}

	configdir, err := utils.ResolveDir(dirpath)
	if err != nil {
		return nil, err
	}

	cfgs, err := configmanager.NewManager(configdir).GetConfigurations()
	if err != nil {
		return nil, fmt.Errorf("couldn't get saved configurations: %v", err)
	}

	if len(cfgs) == 0 {
		return nil, nil
	}

	return cfgs, nil
}

func ListConfigurationTask(ctx context.Context, req *pb.ListConfigurationsRequest) (*pb.ListConfigurationsResponse, error) {
	cfgs, err := getConfigs()
	if err != nil {
		return &pb.ListConfigurationsResponse{Result: "\nError while reading configurations found\n"}, nil
	}

	if len(cfgs) == 0 {
		return &pb.ListConfigurationsResponse{Result: "\nNo configurations found\n"}, nil
	}

	// The user can filter for certain entries using Fuzzy matching.
	searchPattern := req.SearchPattern

	if searchPattern != "" {
		cfgs = configmanager.Entries(cfgs).Filter(func(c *configmanager.Entry) bool {
			return fuzzy.Match(strings.ToLower(searchPattern), strings.ToLower(c.Name))
		})
	}

	if len(cfgs) == 0 {
		return &pb.ListConfigurationsResponse{
			Result: fmt.Sprintf("\nNo configurations found with search pattern \"%s\" \n", searchPattern),
		}, nil
	}

	var output strings.Builder

	output.WriteString("\n")

	for i := range cfgs {

		if i != 0 {
			output.WriteString("\n")
		}

		// config is prented without a new line at its end.
		writeConfigToOutput(&output, cfgs[i])
		output.WriteString("\n")
	}

	return &pb.ListConfigurationsResponse{Result: output.String()}, nil
}

func writeConfigToOutput(out *strings.Builder, entry configmanager.Entry) {
	template := `%s
  - SSH server:  		%s
  - User:        		%s
  - Private key: 		%s
  - Remote:      		%s:%d
  - Default Local Port:		%s`
	nameAndDesc := utils.Bold(entry.Name)
	if strings.TrimSpace(entry.Description) != "" {
		nameAndDesc += " " + "(" + entry.Description + ")"
	}
	out.Write([]byte(
		fmt.Sprintf(
			template,
			nameAndDesc,
			entry.Server,
			entry.User,
			entry.KeyFile,
			entry.RemoteHost,
			entry.RemotePort,
			utils.IntToString(entry.LocalPort),
		)))
}