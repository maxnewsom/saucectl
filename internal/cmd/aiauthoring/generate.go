package aiauthoring

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	cmds "github.com/saucelabs/saucectl/internal/cmd"
	"github.com/saucelabs/saucectl/internal/aiauthoring"
	"github.com/saucelabs/saucectl/internal/usage"
	"github.com/spf13/cobra"
)

func GenerateCommand() *cobra.Command {
	var name string
	var intent string
	var testURL string
	var browserName string
	var platformName string
	var browserVersion string
	var capabilitiesJSON string
	var maxSteps int
	var timeout time.Duration
	var testSuiteID string
	var wait bool
	var pollInterval time.Duration

	cmd := &cobra.Command{
		Use:          "generate",
		Short:        "Generate an AI-authored test case from a natural-language prompt",
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, _ []string) {
			tracker := usage.DefaultClient
			go func() {
				tracker.Collect(cmds.FullName(cmd), usage.Flags(cmd.Flags()))
				_ = tracker.Close()
			}()
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			if name == "" {
				return errors.New("--name is required")
			}
			if intent == "" {
				return errors.New("--intent is required")
			}

			caps := map[string]interface{}{
				"browserName":    browserName,
				"platformName":   platformName,
				"browserVersion": browserVersion,
			}
			if capabilitiesJSON != "" {
				caps = map[string]interface{}{}
				if err := json.Unmarshal([]byte(capabilitiesJSON), &caps); err != nil {
					return fmt.Errorf("failed to parse --capabilities: %w", err)
				}
			}

			req := aiauthoring.GenerateRequest{
				Name: name,
				RunSettings: aiauthoring.RunSettings{
					Target:  aiauthoring.Target{Capabilities: caps},
					TestURL: testURL,
				},
				PromptSettings: aiauthoring.PromptSettings{
					Intent:   intent,
					MaxSteps: maxSteps,
				},
				TestSuiteID: testSuiteID,
			}
			if timeout > 0 {
				req.Timeout = int(timeout.Milliseconds())
			}

			ctx := cmd.Context()
			taskID, err := client.GenerateTestCase(ctx, req)
			if err != nil {
				return fmt.Errorf("failed to start generation: %w", err)
			}
			fmt.Printf("Generation started. Task ID: %s\n", taskID)

			if !wait {
				return nil
			}

			for {
				status, err := client.GetGenerationStatus(ctx, taskID)
				if err != nil {
					return fmt.Errorf("failed to get generation status: %w", err)
				}

				switch status.Data.Status {
				case aiauthoring.StatusCompleted:
					fmt.Printf("Generation complete. Test case ID: %s\n", status.Data.TestCaseID)
					return nil
				case aiauthoring.StatusFailed:
					detail := "unknown error"
					if status.Data.Error != nil {
						detail = status.Data.Error.Detail
					}
					return fmt.Errorf("generation failed: %s", detail)
				default:
					fmt.Printf("Status: %s...\n", status.Data.Status)
					time.Sleep(pollInterval)
				}
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&name, "name", "", "Test case name (required)")
	flags.StringVar(&intent, "intent", "", "Plain-language description of what the test should do (required)")
	flags.StringVar(&testURL, "test-url", "", "Starting URL for a web test")
	flags.StringVar(&browserName, "browser", "chrome", "Browser name (web tests)")
	flags.StringVar(&platformName, "platform", "Windows 11", "Platform name")
	flags.StringVar(&browserVersion, "browser-version", "latest", "Browser version")
	flags.StringVar(&capabilitiesJSON, "capabilities", "", "Raw JSON W3C capabilities, overrides --browser/--platform/--browser-version (needed for mobile)")
	flags.IntVar(&maxSteps, "max-steps", 0, "Cap on AI actions (default 50, max 200)")
	flags.DurationVar(&timeout, "timeout", 0, "Generation timeout (default 5m, max 1h)")
	flags.StringVar(&testSuiteID, "test-suite-id", "", "Attach the generated test case to an existing suite")
	flags.BoolVar(&wait, "wait", true, "Wait for generation to complete before returning")
	flags.DurationVar(&pollInterval, "poll-interval", 5*time.Second, "Polling interval while waiting")

	return cmd
}
