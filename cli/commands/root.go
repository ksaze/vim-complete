package commands

import (
	"github.com/ksaze/vim-complete/internal/toolchain"
	"github.com/spf13/cobra"
)

var root toolchain.Root

var rootCmd = &cobra.Command{
	Use:          "vc",
	Short:        "Vim motions-based esolang for the future.",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		r, err := toolchain.Locate()
		if err != nil {
			return err
		}
		root = r
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}
