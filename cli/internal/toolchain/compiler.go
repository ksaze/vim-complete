package toolchain

import "path/filepath"

func Compile(root Root, src, objOut string) error {
	src, err := filepath.Abs(src)
	if err != nil {
		return err
	}
	objOut, err = filepath.Abs(objOut)
	if err != nil {
		return err
	}
	return Run(root, root.CompilerScript(), src, objOut)
}
