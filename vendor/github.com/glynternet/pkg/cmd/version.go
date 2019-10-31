package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

// NewVersionCmd provides a cobra Command that will print out the version provided as an argument
// The command can be added to a parent command to easily create a version subcommand
func NewVersionCmd(version string, w io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "show the version of this application",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := w.Write([]byte(version))
			return err
		},
	}
}
