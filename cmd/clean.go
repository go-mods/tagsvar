package cmd

import (
	"github.com/go-mods/tagsvar/modules/config"
	"github.com/go-mods/tagsvar/modules/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

// clean command options
type cleanOptions struct {
	Dir         string
	IsRecursive bool
}

// clean command
func newCleanCmd() *cobra.Command {

	o := &cleanOptions{}

	cleanCmd := &cobra.Command{
		Use:     "clean",
		Aliases: []string{"c"},
		Short:   "delete generated files",
		PreRun: func(cmd *cobra.Command, args []string) {
			if config.C.Silent {
				log.Logger = log.Logger.Level(zerolog.Disabled)
			} else if config.C.Verbose {
				log.Logger = log.Logger.Level(zerolog.DebugLevel)
			}
		},
		Run: o.clean,
	}

	// Add flags
	cleanCmd.Flags().StringVarP(&o.Dir, "dir", "d", ".", "Delete generated files in the directory")
	cleanCmd.Flags().BoolVarP(&o.IsRecursive, "recursive", "r", false, "Recursively delete generated files in all subdirectories")
	cleanCmd.Flags().BoolVarP(&config.C.Verbose, "verbose", "v", false, "Print files being deleted")
	cleanCmd.Flags().BoolVarP(&config.C.Silent, "silent", "s", false, "Do not print anything")

	return cleanCmd
}

// clean command
func (o *cleanOptions) clean(cmd *cobra.Command, args []string) {
	var err error

	// Get the working directory
	o.Dir, err = fs.WorkDir(o.Dir)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not get working directory")
		return
	}

	// Info message
	log.Info().Msgf("Cleaning directory %s", o.Dir)

	// List files to delete
	files, err := fs.ListFiles(o.Dir, o.IsRecursive, config.C.ValidateFileName)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not list files")
		return
	}

	// Delete files
	for _, file := range files {
		log.Debug().Msgf("Deleting file : %s", file)
		err = os.Remove(file)
		if err != nil {
			log.Fatal().Err(err).Msgf("Could not delete file %s", file)
			return
		}
	}

	// Info message
	log.Info().Msgf("Finished cleaning directory %s", o.Dir)
}
