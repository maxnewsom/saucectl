package aiauthoring

import (
	"errors"
	"time"

	"github.com/saucelabs/saucectl/internal/credentials"
	"github.com/saucelabs/saucectl/internal/http"
	"github.com/saucelabs/saucectl/internal/region"
	"github.com/saucelabs/saucectl/internal/usage"
	"github.com/spf13/cobra"
)

var client http.AIAuthoring

// Command creates the `ai-authoring` command tree.
func Command(preRun func(cmd *cobra.Command, args []string)) *cobra.Command {
	var regio string

	cmd := &cobra.Command{
		Use:              "ai-authoring",
		Short:            "Interact with Sauce Labs AI Test Authoring",
		SilenceUsage:     true,
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if preRun != nil {
				preRun(cmd, args)
			}

			reg := region.FromString(regio)
			if reg == region.None {
				return errors.New("invalid region")
			}
			if reg == region.Staging {
				usage.DefaultClient.Enabled = false
			}

			creds := credentials.Get()
			client = http.NewAIAuthoring(reg.APIBaseURL(), creds.Username, creds.AccessKey, 15*time.Minute)

			return nil
		},
	}

	flags := cmd.PersistentFlags()
	flags.StringVarP(&regio, "region", "r", "us-west-1", "The Sauce Labs region. Options: us-west-1, us-east-4, eu-central-1.")

	cmd.AddCommand(
		GenerateCommand(),
		ListCommand(),
		GetCommand(),
		RunCommand(),
		PullCommand(),
	)

	return cmd
}
