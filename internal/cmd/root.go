package cmd

import (
	"github.com/0x00-ketsu/taskcli/internal/cmd/flags"
	"github.com/0x00-ketsu/taskcli/internal/layout"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

// rootCmd represents the root command
var rootCmd = &cobra.Command{
	Use:   "taskcli",
	Short: "A terminal UI for manage tasks",
	Run: func(cmd *cobra.Command, args []string) {
		// Load layout
		app := tview.NewApplication()
		layout := layout.Load(app)
		if err := app.SetRoot(layout, true).Run(); err != nil {
			panic(err)
		}
	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flags.Storage, "storage", "s", "~/.taskcli/bolt.db", "taskcli data storage location")
	rootCmd.PersistentFlags().StringVarP(&flags.Editor, "editor", "c", "vim", "external editor for task detail panel")
}
