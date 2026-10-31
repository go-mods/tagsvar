// Package cmd holds the command line interface of tagsvar: the root command and its gen and clean sub-commands.
package cmd

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/go-mods/tagsvar/modules/config"
	"github.com/go-mods/tagsvar/modules/logger"
)

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tagsvar",
		Short: "Generate Go constants and variables from struct tags",
		Long: `tagsvar reads the tags of the structs marked with a "// #tagsvar" comment
and generates, next to each source file, a variables file holding a constant
for each tag name and a map for the options of each tag.`,
		Version: config.C.Version,
		// Errors are printed once by main, without the usage
		SilenceUsage:  true,
		SilenceErrors: true,
		// Runs after the flags of the sub-command are parsed
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			level := zerolog.InfoLevel
			if config.C.Silent {
				level = zerolog.Disabled
			} else if config.C.Verbose {
				level = zerolog.DebugLevel
			}
			log.Logger = logger.New(cmd.ErrOrStderr()).Level(level)
		},
	}

	// Add sub-commands
	rootCmd.AddCommand(newCleanCmd())
	rootCmd.AddCommand(newGenCmd())

	// No completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd
}

// Execute runs the command given on the command line and returns its error.
// The error is not printed: the caller prints it and sets the exit code.
func Execute() error {
	return newRootCmd().Execute()
}

// addOutputFlags adds the flags which set the level of the messages
func addOutputFlags(cmd *cobra.Command, verboseUsage string) {
	// The defaults come from the configuration, so that TAGSVAR_VERBOSE and TAGSVAR_SILENT apply
	cmd.Flags().BoolVarP(&config.C.Verbose, "verbose", "v", config.C.Verbose, verboseUsage)
	cmd.Flags().BoolVarP(&config.C.Silent, "silent", "s", config.C.Silent, "Do not print anything, except errors")
}
