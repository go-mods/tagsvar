package cmd

import (
	"github.com/go-mods/tagsvar/modules/config"
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "tagsvar",
		Short:   "todo",
		Long:    "todo",
		Version: config.C.Version,
	}

	// Add sub-commands
	rootCmd.AddCommand(newCleanCmd())
	rootCmd.AddCommand(newGenCmd())

	//
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd
}

func Execute() (err error) {
	rootCmd := newRootCmd()
	if err = rootCmd.Execute(); err != nil {
		return
	}
	return nil
}
