package cmd

import (
	"context"

	"github.com/render-oss/cli/pkg/client"
	"github.com/render-oss/cli/pkg/service"
	"github.com/render-oss/cli/pkg/tui"
	"github.com/render-oss/cli/pkg/tui/views"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/render-oss/cli/pkg/command"
)

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh [serviceID]",
	Short: "SSH into a service instance",
	Long: `SSH into a service instance. Optionally pass the service id as an argument.
To pass arguments to ssh, use the following syntax: render ssh [serviceID] -- [ssh args]`,
	GroupID: GroupSession.ID,
}

func InteractiveSSHView(ctx context.Context, input *views.SSHInput, breadcrumb string) tea.Cmd {
	return command.AddToStackFunc(
		ctx,
		sshCmd,
		breadcrumb,
		input,
		views.NewSSHView(ctx, input, tui.WithCustomOptions[*service.Model](getSSHTableOptions(ctx, breadcrumb))),
	)
}

func getSSHTableOptions(ctx context.Context, breadcrumb string) []tui.CustomOption {
	return []tui.CustomOption{
		WithCopyID(ctx, servicesCmd),
		WithWorkspaceSelection(ctx),
		WithProjectFilter(ctx, servicesCmd, "Project Filter", &views.SSHInput{}, func(ctx context.Context, project *client.Project) tea.Cmd {
			input := views.SSHInput{}
			if project != nil {
				input.Project = project
				input.EnvironmentIDs = project.EnvironmentIds
			}
			return InteractiveSSHView(ctx, &input, breadcrumb)
		}),
	}
}

func init() {
	rootCmd.AddCommand(sshCmd)

	sshCmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		input := views.SSHInput{}
		err := command.ParseCommandInteractiveOnly(cmd, args, &input)
		if err != nil {
			return err
		}

		if cmd.ArgsLenAtDash() == 0 {
			input.ServiceID = ""
		}

		if cmd.ArgsLenAtDash() >= 0 {
			input.Args = args[cmd.ArgsLenAtDash():]
		}

		InteractiveSSHView(ctx, &input, "SSH")
		return nil
	}
}
