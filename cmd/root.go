package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "secondseat",
	Short:        "Add or remove input devices for a second seat.",
	Long:         "Add or remove input devices for a second seat.",
	SilenceUsage: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.CompletionOptions.DisableNoDescFlag = true
	rootCmd.AddCommand(newAddCmd())
}
