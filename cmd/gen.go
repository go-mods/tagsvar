package cmd

import (
	"github.com/go-mods/tagsvar/modules/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// gen command options
type genOptions struct {
	Dir         string
	IsRecursive bool
}

// clean command
func newGenCmd() *cobra.Command {

	o := &genOptions{}

	genCmd := &cobra.Command{
		Use:     "gen",
		Aliases: []string{"g"},
		Short:   "todo",
		Long:    "todo",
		PreRun: func(cmd *cobra.Command, args []string) {
			if config.C.Silent {
				log.Logger = log.Logger.Level(zerolog.Disabled)
			} else if config.C.Verbose {
				log.Logger = log.Logger.Level(zerolog.DebugLevel)
			}
		},
		Run: o.gen,
	}

	// Add flags
	genCmd.Flags().StringVarP(&o.Dir, "dir", "d", ".", "Generate variables files for the directory")
	genCmd.Flags().BoolVarP(&o.IsRecursive, "recursive", "r", false, "Generate variables files for all subdirectories")
	genCmd.Flags().BoolVarP(&config.C.Verbose, "verbose", "v", false, "Print files being deleted")
	genCmd.Flags().BoolVarP(&config.C.Silent, "silent", "s", false, "Do not print anything")

	return genCmd
}

// gen command
func (o *genOptions) gen(cmd *cobra.Command, args []string) {
}
