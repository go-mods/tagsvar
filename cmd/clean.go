package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/go-mods/tagsvar/modules/fs"
)

// clean command options
type cleanOptions struct {
	Dir         string
	IsRecursive bool
}

// newCleanCmd creates the clean command
func newCleanCmd() *cobra.Command {
	o := &cleanOptions{}

	cleanCmd := &cobra.Command{
		Use:     "clean",
		Aliases: []string{"c"},
		Short:   "Delete the generated variables files of a directory",
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return o.clean()
		},
	}

	// Add flags
	cleanCmd.Flags().StringVarP(&o.Dir, "dir", "d", ".", "Delete generated files in the directory")
	cleanCmd.Flags().BoolVarP(&o.IsRecursive, "recursive", "r", false, "Recursively delete generated files in all subdirectories")
	addOutputFlags(cleanCmd, "Print files being deleted")

	return cleanCmd
}

// clean deletes the generated files of the directory
func (o *cleanOptions) clean() error {
	dir, err := fs.WorkDir(o.Dir)
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	log.Info().Msgf("Cleaning directory %s", dir)

	files, err := fs.ListFiles(dir, o.IsRecursive, fs.IsGeneratedFile)
	if err != nil {
		return fmt.Errorf("list generated files in %s: %w", dir, err)
	}

	for _, file := range files {
		log.Debug().Msgf("Deleting file %s", file)
		if err := os.Remove(file); err != nil {
			return fmt.Errorf("delete generated file: %w", err)
		}
	}

	log.Info().Msgf("Finished cleaning directory %s", dir)
	return nil
}
