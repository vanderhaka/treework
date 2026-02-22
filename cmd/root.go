package cmd

import (
	"fmt"
	"os"

	"github.com/jamesvanderhaak/wt/internal/config"
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
	fmt.Println()
	fmt.Println(ui.Banner())

	// First-run setup: prompt for base folder if no config file and no DEV_DIR
	if !config.FileExists() && os.Getenv("DEV_DIR") == "" {
		fmt.Println()
		ui.Info("Welcome! Let's set your base folder (where your git repos live).")
		fmt.Println()

		home, _ := os.UserHomeDir()
		selected, err := SetBaseDir(home)
		if err != nil {
			// User escaped — save defaults so we don't prompt again
			_ = config.Save(&config.Config{})
		} else {
			_ = config.Save(&config.Config{BaseDir: selected})
			ui.Success(fmt.Sprintf("Base folder set to %s", selected))
		}
	}

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
		case "settings":
			runSettingsInteractive()
		case "quit":
			fmt.Println()
			ui.Muted("Goodbye.")
			fmt.Println()
			return
		}
	}
}
