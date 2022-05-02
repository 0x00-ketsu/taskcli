/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/0x00-ketsu/taskcli/cmd/flags"
	"github.com/0x00-ketsu/taskcli/internal/database"
	"github.com/0x00-ketsu/taskcli/internal/global"
	"github.com/0x00-ketsu/taskcli/internal/layout"
	"github.com/0x00-ketsu/taskcli/internal/repository/bolt"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

// rootCmd represents the root command
var rootCmd = &cobra.Command{
	Use:   "taskcli",
	Short: "A terminal UI for manage tasks",
	Run: func(cmd *cobra.Command, args []string) {
		// Create app
		app := tview.NewApplication()

		// Initial global variables
		global.App = app

		// Connect DB
		db := database.Connect(flags.Storage)
		defer func() {
			if err := db.Close(); err != nil {
				panic(err)
			}
		}()

		// Task repository
		taskRepo := bolt.NewTask(db)
		global.TaskRepo = taskRepo

		// Load layout
		layout := layout.Load()
		global.Layout = layout

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
