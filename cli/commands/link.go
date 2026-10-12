package commands

import (
	"fmt"
	"os"

	"github.com/ksaze/vim-complete/internal/discover"
	"github.com/ksaze/vim-complete/internal/paths"
	"github.com/ksaze/vim-complete/internal/stale"
	"github.com/ksaze/vim-complete/internal/toolchain"
	"github.com/spf13/cobra"
)

var linkOut string

func init() {
	linkCmd.Flags().StringVarP(&linkOut, "out", "o", "a.vcx", "output executable name")
	rootCmd.AddCommand(linkCmd)
}

var linkCmd = &cobra.Command{
	Use:   "link <main-file>",
	Short: "Link the main module against all compiled objects",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mainFile := args[0]
		if err := paths.MustExist(mainFile); err != nil {
			return err
		}
		if err := paths.EnsureBuildDir(); err != nil {
			return err
		}

		localObjs, err := discover.Files(paths.BuildDir, paths.ObjExt)
		if err != nil {
			return err
		}
		if len(localObjs) == 0 {
			return fmt.Errorf("no object files found in %s (run compile first)", paths.BuildDir)
		}
		stdlibObjs, err := discover.Files(root.StdlibDir, paths.ObjExt)
		if err != nil {
			return err
		}

		allObjs := append(append([]string{mainFile}, localObjs...), stdlibObjs...)
		if !stale.IsStale(linkOut, allObjs...) {
			fmt.Println("up to date:", linkOut)
			return nil
		}

		fmt.Println("linking ->", linkOut)
		if err := toolchain.Link(root, mainFile, linkOut); err != nil {
			return err
		}
		if _, err := os.Stat(linkOut); err != nil {
			return fmt.Errorf("linker reported success but did not produce %s (check linker.vcx)", linkOut)
		}
		return nil
	},
}
