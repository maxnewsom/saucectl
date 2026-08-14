package aiauthoring

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	cmds "github.com/saucelabs/saucectl/internal/cmd"
	"github.com/saucelabs/saucectl/internal/tables"
	"github.com/saucelabs/saucectl/internal/usage"
	"github.com/spf13/cobra"
)

func ListCommand() *cobra.Command {
	var search string
	var limit int
	var skip int

	cmd := &cobra.Command{
		Use:          "list",
		Short:        "List AI-authored test cases",
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, _ []string) {
			tracker := usage.DefaultClient
			go func() {
				tracker.Collect(cmds.FullName(cmd), usage.Flags(cmd.Flags()))
				_ = tracker.Close()
			}()
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			testCases, err := client.ListTestCases(cmd.Context(), search, limit, skip)
			if err != nil {
				return fmt.Errorf("failed to list test cases: %w", err)
			}

			t := table.NewWriter()
			t.SetStyle(tables.DefaultTableStyle)
			t.SuppressEmptyColumns()
			t.AppendHeader(table.Row{"ID", "Name", "Created"})
			for _, tc := range testCases {
				t.AppendRow(table.Row{tc.ID, tc.Name, tc.CreationDate})
			}
			fmt.Println(t.Render())
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&search, "search", "", "Filter by name")
	flags.IntVar(&limit, "limit", 20, "Max results to return")
	flags.IntVar(&skip, "skip", 0, "Number of results to skip")

	return cmd
}
