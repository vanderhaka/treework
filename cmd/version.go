package cmd

import (
	"fmt"

	"github.com/jamesvanderhaak/wt/internal/ui"
	"github.com/spf13/cobra"
)

var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(ui.BrandStyle.Render("treework") + " " + ui.MutedStyle.Render("v"+Version))
	},
}
