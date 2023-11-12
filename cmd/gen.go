package cmd

import "github.com/spf13/cobra"

// gen command options
type genOptions struct {
}

// clean command
func newGenCmd() *cobra.Command {

	o := &genOptions{}

	genCmd := &cobra.Command{
		Use:     "gen",
		Aliases: []string{"g"},
		Short:   "todo",
		Long:    "todo",
		RunE:    o.gen,
	}

	return genCmd
}

// gen command
func (o *genOptions) gen(cmd *cobra.Command, args []string) error {

	return nil
}
