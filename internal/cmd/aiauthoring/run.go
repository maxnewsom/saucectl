package aiauthoring

import (
	"errors"
	"fmt"

	"github.com/saucelabs/saucectl/internal/aiauthoring"
	cmds "github.com/saucelabs/saucectl/internal/cmd"
	"github.com/saucelabs/saucectl/internal/usage"
	"github.com/spf13/cobra"
)

func RunCommand() *cobra.Command {
	var buildName string
	var scTunnelName string

	cmd := &cobra.Command{
		Use:          "run <testCaseID>",
		Short:        "Run a saved AI-authored test case",
		SilenceUsage: true,
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 || args[0] == "" {
				return errors.New("no test case ID specified")
			}
			return nil
		},
		PreRun: func(cmd *cobra.Command, _ []string) {
			tracker := usage.DefaultClient
			go func() {
				tracker.Collect(cmds.FullName(cmd), usage.Flags(cmd.Flags()))
				_ = tracker.Close()
			}()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := aiauthoring.RunRequest{
				BuildName:    buildName,
				SCTunnelName: scTunnelName,
			}

			resp, err := client.RunTestCase(cmd.Context(), args[0], req)
			if err != nil {
				return fmt.Errorf("failed to run test case: %w", err)
			}

			fmt.Printf("Run started. Run ID: %s\n", resp.Data.ID)
			for _, job := range resp.Data.Jobs {
				fmt.Printf("  Job: %s\n", job.Name)
				if job.URL != "" {
					fmt.Printf("    %s\n", job.URL)
				}
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&buildName, "build-name", "", "Sauce Labs build name to group the run under")
	flags.StringVar(&scTunnelName, "sc-tunnel-name", "", "Name of an active Sauce Connect tunnel")

	return cmd
}
