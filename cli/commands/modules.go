package commands

import (
	"fmt"

	"github.com/ksaze/vim-complete/internal/discover"
	"github.com/ksaze/vim-complete/internal/paths"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(modulesCmd)
}

var modulesCmd = &cobra.Command{
	Use:   "modules",
	Short: "List local and stdlib modules visible to the linker",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		localObjs, err := discover.Files(paths.BuildDir, paths.ObjExt)
		if err != nil {
			return fmt.Errorf("discover local objects: %w", err)
		}
		stdlibObjs, err := discover.Files(root.StdlibDir, paths.ObjExt)
		if err != nil {
			return fmt.Errorf("discover stdlib objects: %w", err)
		}

		if len(localObjs) == 0 && len(stdlibObjs) == 0 {
			fmt.Println("no modules found")
			return nil
		}

		for _, p := range localObjs {
			fmt.Println("local:", p)
		}
		for _, p := range stdlibObjs {
			fmt.Println("std:", p)
		}
		return nil
	},
}
