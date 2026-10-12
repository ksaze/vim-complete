package commands

import (
	"path/filepath"

	"github.com/ksaze/vim-complete/internal/paths"
	"github.com/ksaze/vim-complete/internal/toolchain"
	"github.com/spf13/cobra"
)

func init() { rootCmd.AddCommand(debugCmd) }

var debugCmd = &cobra.Command{
	Use:          "debug <target> [args...]",
	Short:        "Run a compiled executable under the interactive debugger",
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		target, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}
		if err := paths.MustExist(target); err != nil {
			return err
		}
		return toolchain.RunInteractive(root, target, args[1:]...)
	},
}
