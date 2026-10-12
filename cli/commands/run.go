package commands

import (
	"github.com/ksaze/vim-complete/internal/paths"
	"github.com/ksaze/vim-complete/internal/toolchain"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run <target> [args...]",
	Short: "Run a compiled executable on the runtime",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		if err := paths.MustExist(target); err != nil {
			return err
		}
		return toolchain.Run(root, target, args[1:]...)
	},
}
