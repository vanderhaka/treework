package cmd

import (
	"fmt"
	"os"

	"github.com/jamesvanderhaak/wt/internal/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "treework",
	Short: "Git worktree manager",
	Long:  ui.Banner(),
	Run:   runRoot,
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(clearCmd)
	rootCmd.AddCommand(versionCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runRoot(cmd *cobra.Command, args []string) {
	for {
		fmt.Println()
		action, err := ui.SelectAction()
		if err != nil {
			handleAbort(err)
			return
		}

		switch action {
		case "new":
			runNewInteractive(cmd)
		case "ls":
			runLsInteractive(cmd)
		case "rm":
			runRmInteractive(cmd)
		case "clear":
			runClearInteractive(cmd)
		case "quit":
			fmt.Println()
			ui.Muted("Goodbye.")
			fmt.Println()
			return
		}
	}
}
