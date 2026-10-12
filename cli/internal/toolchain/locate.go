package toolchain

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Root is a resolved install: everything the toolchain needs, located once.
type Root struct {
	ToolchainDir string // .../toolchain (runtime.lua + compiler.vcx + linker.vcx)
	StdlibDir    string // .../stdlib
}

const scriptExt = ".vcx"

func (r Root) RuntimeLua() string     { return filepath.Join(r.ToolchainDir, "runtime.lua") }
func (r Root) DebuggerLua() string    { return filepath.Join(r.ToolchainDir, "debugger.lua") }
func (r Root) CompilerScript() string { return filepath.Join(r.ToolchainDir, "compiler"+scriptExt) }
func (r Root) LinkerScript() string   { return filepath.Join(r.ToolchainDir, "linker"+scriptExt) }

// NvimBin resolves the nvim binary used to drive the runtime.
// NVIM_BIN overrides (pin a specific build/version); otherwise PATH.
func NvimBin() (string, error) {
	if p := os.Getenv("NVIM_BIN"); p != "" {
		return p, nil
	}
	p, err := exec.LookPath("nvim")
	if err != nil {
		return "", fmt.Errorf("nvim not found on PATH (required as the vc runtime): %w", err)
	}
	return p, nil
}

// Locate resolves the toolchain/stdlib root relative to the running vc binary.
// VC_ROOT overrides this entirely.
func Locate() (Root, error) {
	if envRoot := os.Getenv("VC_ROOT"); envRoot != "" {
		root := rootFrom(envRoot)
		if err := root.validate(); err != nil {
			return Root{}, fmt.Errorf("VC_ROOT=%s: %w", envRoot, err)
		}
		return root, nil
	}

	exe, err := os.Executable()
	if err != nil {
		return Root{}, fmt.Errorf("locate toolchain: %w", err)
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return Root{}, fmt.Errorf("locate toolchain: %w", err)
	}

	root := rootFrom(filepath.Dir(exe))
	if err := root.validate(); err != nil {
		return Root{}, fmt.Errorf("toolchain not found next to %s: %w (set VC_ROOT to override)", exe, err)
	}
	return root, nil
}

func rootFrom(base string) Root {
	return Root{
		ToolchainDir: filepath.Join(base, "toolchain"),
		StdlibDir:    filepath.Join(base, "stdlib"),
	}
}

func (r Root) validate() error {
	for _, p := range []string{r.ToolchainDir, r.StdlibDir, r.RuntimeLua(), r.CompilerScript(), r.LinkerScript()} {
		if _, err := os.Stat(p); err != nil {
			return fmt.Errorf("missing %s", p)
		}
	}
	return nil
}
