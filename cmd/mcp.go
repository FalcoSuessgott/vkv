package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

type mcpImpl struct {
	binary string
}

func NewMCPCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "mcp",
		Short:         "start a MCP server that provides vkv capabilities",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			m := mcpImpl{}

			bin, err := os.Executable()
			if err != nil {
				return err
			}

			m.binary = bin

			s := server.NewMCPServer("vkv", Version)

			s.AddTool(
				mcp.NewTool("export",
					mcp.WithDescription("Export secrets from Vault KV engine"),
					mcp.WithString("enginePath",
						mcp.Required(),
						mcp.Description("Path to the KV engine"),
					),
				),
				m.export,
			)

			return server.ServeStdio(s)
		},
	}
}

func (m *mcpImpl) export(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	enginePath, err := req.RequireString("enginePath")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	out, err := exec.CommandContext(ctx, m.binary, "export", "-f=markdown", "--show-values", fmt.Sprintf("-p=%s", enginePath)).CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(out)), nil
}
