package toolchain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const runtimeTimeout = 30 * time.Second

type scriptResult struct {
	ExitCode int    `json:"exit_code"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
}

func Run(root Root, script string, args ...string) error {
	nvim, err := NvimBin()
	if err != nil {
		return err
	}

	captureDir, err := os.MkdirTemp("", "vc-run-*")
	if err != nil {
		return fmt.Errorf("runtime: create capture dir: %w", err)
	}
	defer os.RemoveAll(captureDir)
	captureDir, err = filepath.EvalSymlinks(captureDir)
	if err != nil {
		return fmt.Errorf("runtime: resolve capture dir: %w", err)
	}

	resultPath := filepath.Join(captureDir, "result.json")

	ctx, cancel := context.WithTimeout(context.Background(), runtimeTimeout)
	defer cancel()

	cmdArgs := append([]string{"--headless", "-u", root.RuntimeLua(), "-S", script}, args...)
	cmd := exec.CommandContext(ctx, nvim, cmdArgs...)
	// cmd.Dir intentionally unset — inherits the CLI's real invocation directory.
	cmd.Env = append(os.Environ(), "VC_RESULT="+resultPath)
	runErr := cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("runtime timed out after %s on %s", runtimeTimeout, script)
	}
	if runErr != nil {
		// nvim itself crashed/was killed — no result.json to trust.
		return fmt.Errorf("runtime process failed on %s: %w", script, runErr)
	}

	data, err := os.ReadFile(resultPath)
	if err != nil {
		return fmt.Errorf("runtime did not produce a result for %s (check the script writes VC_RESULT): %w", script, err)
	}

	var res scriptResult
	if err := json.Unmarshal(data, &res); err != nil {
		return fmt.Errorf("runtime result for %s is not valid JSON: %w", script, err)
	}

	if res.Stdout != "" {
		fmt.Fprint(os.Stdout, res.Stdout)
	}
	if res.Stderr != "" {
		fmt.Fprint(os.Stderr, res.Stderr)
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("%s failed (exit %d)", filepath.Base(script), res.ExitCode)
	}
	return nil
}

// RunInteractive runs a compiled target under the debugger with nvim's real
// TUI attached to the user's terminal. Unlike Run, stdio is inherited, there
// is no timeout, and no result.json is read; the session's outcome is nvim's
// own exit status.
func RunInteractive(root Root, target string, args ...string) error {
	nvim, err := NvimBin()
	if err != nil {
		return err
	}
	if _, err := os.Stat(root.DebuggerLua()); err != nil {
		return fmt.Errorf("debugger not available in this install: %w", err)
	}

	cmdArgs := append([]string{"-u", root.DebuggerLua(), "-S", target}, args...)
	cmd := exec.Command(nvim, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return fmt.Errorf("debugger exited with status %d", exitErr.ExitCode())
		}
		return fmt.Errorf("debugger failed: %w", err)
	}
	return nil
}
