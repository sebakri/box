//go:build darwin
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func applySandbox(cmd *exec.Cmd, name string, args []string, rootDir string) (string, []string) {
	parentDir := filepath.Dir(rootDir)
	profile := fmt.Sprintf(`(version 1)
(allow default)
(deny file-write*)
(allow file-write* (subpath %q))
(allow file-write* (subpath %q))
(allow file-write* (subpath %q))
`, rootDir, parentDir, os.TempDir())

	return "sandbox-exec", append([]string{"-p", profile, name}, args...)
}

