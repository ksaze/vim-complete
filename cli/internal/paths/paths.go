package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	SrcExt = ".vc"
	ObjExt = ".vco"
	ExeExt = "vcx"

	BuildDir = "vc-build"
)

// ObjPathIn derives the object path for src within an arbitrary directory,
// e.g. ObjPathIn("mylibs", "foo.vl") -> mylibs/foo.vco.
func ObjPathIn(dir, src string) string {
	return filepath.Join(dir, replaceExt(filepath.Base(src), ObjExt))
}

// ObjPath derives foo.yl -> foo.ylo
func ObjPath(src string) string {
	return filepath.Join(BuildDir, replaceExt(filepath.Base(src), ObjExt))
}

// ExePath derives foo.yl -> foo.ylx (or from an obj path)
func ExePath(src string) string {
	return replaceExt(filepath.Base(src), ExeExt)
}

func replaceExt(p, newExt string) string {
	ext := filepath.Ext(p)
	return strings.TrimSuffix(p, ext) + newExt
}

func EnsureBuildDir() error {
	if err := os.MkdirAll(BuildDir, 0755); err != nil {
		return fmt.Errorf("create %s %w", BuildDir, err)
	}

	return nil
}

// MustExist returns an error if p does not exist or cannot be accessed.
func MustExist(p string) error {
	if _, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no such file: %s", p)
		}
		return fmt.Errorf("cannot access %s: %w", p, err)
	}
	return nil
}
