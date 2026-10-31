package cmd

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/go-mods/tagsvar/modules/fs"
	"github.com/go-mods/tagsvar/modules/generator"
	"github.com/go-mods/tagsvar/modules/parser"
)

// gen command options
type genOptions struct {
	Dir         string
	IsRecursive bool
}

// newGenCmd creates the gen command
func newGenCmd() *cobra.Command {
	o := &genOptions{}

	genCmd := &cobra.Command{
		Use:     "gen",
		Aliases: []string{"g"},
		Short:   "Generate the variables files of a directory",
		Long: `Generate a variables file next to each Go file of the directory which holds
structs marked with a "// #tagsvar" comment. Test files and generated files are skipped.`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return o.gen()
		},
	}

	// Add flags
	genCmd.Flags().StringVarP(&o.Dir, "dir", "d", ".", "Generate variables files for the directory")
	genCmd.Flags().BoolVarP(&o.IsRecursive, "recursive", "r", false, "Generate variables files for all subdirectories")
	addOutputFlags(genCmd, "Print debug messages")

	return genCmd
}

// gen generates the variables files of the directory
func (o *genOptions) gen() error {
	dir, err := fs.WorkDir(o.Dir)
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	log.Info().Msgf("Parsing files in %s", dir)

	parsedFiles, err := parser.NewParser().ParseDir(dir, o.IsRecursive)
	if err != nil {
		return fmt.Errorf("parse files in %s: %w", dir, err)
	}

	if err := generator.NewGenerator().Generate(parsedFiles); err != nil {
		return fmt.Errorf("generate variables files: %w", err)
	}
	return nil
}
