package cmd

import (
	"github.com/go-mods/tagsvar/modules/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

// clean command options
type cleanOptions struct {
	Dir         string
	IsRecursive bool
	Verbose     bool
}

// clean command
func newCleanCmd() *cobra.Command {

	o := &cleanOptions{}

	cleanCmd := &cobra.Command{
		Use:     "clean",
		Aliases: []string{"c"},
		Short:   "delete generated files",
		PreRun: func(cmd *cobra.Command, args []string) {
			if o.Verbose {
				log.Logger = log.Logger.Level(zerolog.DebugLevel)
			}
		},
		RunE: o.clean,
	}

	// Add flags
	cleanCmd.Flags().StringVarP(&o.Dir, "dir", "d", ".", "Delete generated files in the directory")
	cleanCmd.Flags().BoolVarP(&o.IsRecursive, "recursive", "r", false, "Recursively delete generated files in all subdirectories")
	cleanCmd.Flags().BoolVarP(&o.Verbose, "verbose", "v", false, "Print files beeing deleted")

	return cleanCmd
}

// clean command
func (o *cleanOptions) clean(cmd *cobra.Command, args []string) error {
	log.Info().Msg("start cleaning")

	// List files to delete
	// The files must start with the prefix defined in the config
	// The files must end with the suffix defined in the config
	// The files are searched in the current directory or in all subdirectories if recursive is true
	// If verbose is true, the files are printed before deletion

	// Get the directory
	if o.Dir == "" || o.Dir == "." {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal().Err(err).Msg("Could not get current directory")
			return err
		}
		o.Dir = cwd
	} else {
		cwd, err := filepath.Abs(o.Dir)
		if err != nil {
			log.Fatal().Err(err).Msgf("Could not get absolute path of directory %s", o.Dir)
			return err
		}
		o.Dir = cwd
	}
	log.Debug().Msgf("Directory to clean : %s", o.Dir)

	// List files to delete
	err := filepath.Walk(o.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories if recursive is false
		if info.IsDir() && !o.IsRecursive {
			return filepath.SkipDir
		}

		// Get the file name
		fileName := info.Name()

		// Skip files that do not start with the prefix and do not end with the suffix.
		// If the file matches the criteria then delete it.
		if strings.HasPrefix(fileName, config.C.Prefix) && strings.HasSuffix(fileName, config.C.Suffix+".go") {
			log.Debug().Msgf("File to delete : %s", fileName)
			err = os.Remove(path)
			if err != nil {
				return err
			}
			log.Debug().Msgf("File deleted : %s", fileName)
		}
		return nil
	})

	if err != nil {
		log.Fatal().Err(err).Msg("Could not list files")
		return err
	}

	log.Info().Msg("end cleaning")

	return nil
}
