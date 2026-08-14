package aiauthoring

import (
	"errors"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	cmds "github.com/saucelabs/saucectl/internal/cmd"
	"github.com/saucelabs/saucectl/internal/tables"
	"github.com/saucelabs/saucectl/internal/usage"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "get <testCaseID>",
		Short:        "Get an AI-authored test case by ID",
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
			tc, err := client.GetTestCase(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("failed to get test case: %w", err)
			}

			t := table.NewWriter()
			t.SetStyle(tables.DefaultTableStyle)
			t.SuppressEmptyColumns()
			t.AppendHeader(table.Row{"ID", "Name", "Created", "Updated"})
			t.AppendRow(table.Row{tc.ID, tc.Name, tc.CreationDate, tc.LastUpdateDate})
			fmt.Println(t.Render())
			return nil
		},
	}

	return cmd
}
