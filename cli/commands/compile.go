package commands

import (
	"fmt"
	"os"

	"github.com/ksaze/vim-complete/internal/paths"
	"github.com/ksaze/vim-complete/internal/stale"
	"github.com/ksaze/vim-complete/internal/toolchain"
	"github.com/spf13/cobra"
)

var (
	compileOut   string
	compileForce bool
)

func init() {
	compileCmd.Flags().StringVarP(&compileOut, "out", "o", "", "output object file path (default: vc-build/<name>.vco)")
	compileCmd.Flags().BoolVarP(&compileForce, "force", "f", false, "recompile even if the object file is up to date")
	rootCmd.AddCommand(compileCmd)
}

var compileCmd = &cobra.Command{
	Use:   "compile <file>",
	Short: "Compile a source file to an object file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return compileFile(args[0], compileOut, compileForce)
	},
}

// compileFile compiles src to obj (vc-build/<name>.vco if obj is empty).
// force skips the staleness check and always recompiles.
// Shared by `vc compile` and `vc build`.
func compileFile(src, obj string, force bool) error {
	if err := paths.MustExist(src); err != nil {
		return err
	}
	if err := paths.EnsureBuildDir(); err != nil {
		return err
	}

	if obj == "" {
		obj = paths.ObjPath(src)
	}

	if !force && !stale.IsStale(obj, src) {
		fmt.Println("up to date:", obj)
		return nil
	}

	if err := toolchain.Compile(root, src, obj); err != nil {
		return err
	}
	if _, err := os.Stat(obj); err != nil {
		return fmt.Errorf("compiler reported success but did not produce %s (check compiler.vcx)", obj)
	}
	return nil
}
