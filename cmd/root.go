package cmd

import "github.com/spf13/cobra"

var Version = "development"

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "tagsvar",
		Short:   "todo",
		Long:    "todo",
		Version: Version,
	}

	return rootCmd
}

func Execute() (err error) {
	rootCmd := newRootCmd()
	if err = rootCmd.Execute(); err != nil {
		return
	}
	return nil
}
