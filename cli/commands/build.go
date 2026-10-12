package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ksaze/vim-complete/internal/discover"
	"github.com/ksaze/vim-complete/internal/paths"
	"github.com/ksaze/vim-complete/internal/stale"
	"github.com/ksaze/vim-complete/internal/toolchain"
	"github.com/spf13/cobra"
)

var (
	buildMain  string
	buildOut   string
	buildForce bool
)

func init() {
	buildCmd.Flags().StringVarP(&buildMain, "main", "m", "", "entry source file (omit to compile as a library, no link step)")
	buildCmd.Flags().StringVarP(&buildOut, "out", "o", "a.vcx", "executable output path (-m given) or object install directory (library mode)")
	buildCmd.Flags().BoolVarP(&buildForce, "force", "f", false, "recompile/relink even if outputs are up to date")
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build <dir>",
	Short: "Compile every source file in dir, and link if -m is given",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := args[0]

		srcs, err := discover.Files(dir, paths.SrcExt)
		if err != nil {
			return err
		}
		if len(srcs) == 0 {
			return fmt.Errorf("no %s files found in %s", paths.SrcExt, dir)
		}

		// Always compile into vc-build/ first — both modes need every
		// source as an intermediate object before any linking happens,
		// since library files may depend on sibling modules.
		for _, src := range srcs {
			if err := compileFile(src, "", buildForce); err != nil {
				return err
			}
		}

		if buildMain == "" {
			return buildLibrary(srcs)
		}
		return buildExecutable(dir)
	},
}

// buildLibrary links every source file individually against the full set
// of local + stdlib objects, so each library file resolves its own
// dependencies on sibling modules, writing each linked result into buildOut.
func buildLibrary(srcs []string) error {
	if err := os.MkdirAll(buildOut, 0755); err != nil {
		return fmt.Errorf("create %s: %w", buildOut, err)
	}

	localObjs, err := discover.Files(paths.BuildDir, paths.ObjExt)
	if err != nil {
		return err
	}
	stdlibObjs, err := discover.Files(root.StdlibDir, paths.ObjExt)
	if err != nil {
		return err
	}

	for _, src := range srcs {
		obj := paths.ObjPath(src) // vc-build/<name>.vco, from the compile step above
		out := paths.ObjPathIn(buildOut, src)

		allObjs := append(append([]string{obj}, localObjs...), stdlibObjs...)
		if !buildForce && !stale.IsStale(out, allObjs...) {
			fmt.Println("up to date:", out)
			continue
		}

		if err := toolchain.Link(root, obj, out); err != nil {
			return fmt.Errorf("link %s: %w", filepath.Base(src), err)
		}
		if _, err := os.Stat(out); err != nil {
			return fmt.Errorf("linker reported success but did not produce %s (check linker.vcx)", out)
		}
	}

	fmt.Println("built", len(srcs), "module(s) ->", buildOut)
	return nil
}

// buildExecutable links the resolved main object (already compiled into
// vc-build/ by the shared compile step) into buildOut.
func buildExecutable(srcs string) error {
	mainSrc := filepath.Join(srcs, buildMain)
	mainObj := paths.ObjPath(mainSrc)
	if err := paths.MustExist(mainObj); err != nil {
		return fmt.Errorf("main object missing after compile: %w", err)
	}

	localObjs, err := discover.Files(paths.BuildDir, paths.ObjExt)
	if err != nil {
		return err
	}
	stdlibObjs, err := discover.Files(root.StdlibDir, paths.ObjExt)
	if err != nil {
		return err
	}

	allObjs := append(append([]string{mainObj}, localObjs...), stdlibObjs...)
	if !buildForce && !stale.IsStale(buildOut, allObjs...) {
		fmt.Println("up to date:", buildOut)
		return nil
	}

	if err := toolchain.Link(root, mainObj, buildOut); err != nil {
		return err
	}
	if _, err := os.Stat(buildOut); err != nil {
		return fmt.Errorf("linker reported success but did not produce %s (check linker.vcx)", buildOut)
	}
	return nil
}
