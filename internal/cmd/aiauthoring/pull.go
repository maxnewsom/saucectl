package aiauthoring

import (
	"errors"
	"fmt"
	"os"

	cmds "github.com/saucelabs/saucectl/internal/cmd"
	"github.com/saucelabs/saucectl/internal/usage"
	"github.com/spf13/cobra"
)

func PullCommand() *cobra.Command {
	var target string
	var output string

	cmd := &cobra.Command{
		Use:   "pull <testCaseID>",
		Short: "Pull generated source code for a test case onto the local workstation",
		Long: `Pull generated source code for a test case onto the local workstation.

Run without --target to see the valid code-export targets for this test case
(they vary by platform, e.g. javascript_webdriverio for web, java_espresso for
mobile). Then re-run with --target set.`,
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
			id := args[0]
			ctx := cmd.Context()

			if target == "" {
				targets, err := client.GetCodeTargets(ctx, id)
				if err != nil {
					return fmt.Errorf("failed to get code targets: %w", err)
				}
				fmt.Println("Available targets:")
				for _, t := range targets {
					fmt.Printf("  %s\n", t)
				}
				fmt.Println("\nRe-run with --target <name> to pull code.")
				return nil
			}

			code, err := client.GetCode(ctx, id, target)
			if err != nil {
				return fmt.Errorf("failed to get code: %w", err)
			}

			if output == "" {
				fmt.Println(code)
				return nil
			}

			if err := os.WriteFile(output, []byte(code), 0o644); err != nil {
				return fmt.Errorf("failed to write %s: %w", output, err)
			}
			fmt.Printf("Wrote %s\n", output)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&target, "target", "", "Code generation target (e.g. javascript_webdriverio); omit to list valid targets")
	flags.StringVarP(&output, "output", "o", "", "File to write the generated code to (defaults to stdout)")

	return cmd
}
