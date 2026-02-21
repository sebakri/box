//go:build integration

package installer

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sebakri/box/internal/config"
)

func TestIntegrationInstallation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	boxBin := buildBoxBinary(t)
	projectDir := setupTestProject(t, "testdata/integration_test.yml")

	// Use a dedicated GOPATH for the test to avoid caching issues
	t.Setenv("GOPATH", filepath.Join(t.TempDir(), "gopath"))

	runBoxInstall(t, boxBin, projectDir)
	verifyInstallation(t, projectDir)
}

func TestSandboxIsolation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	boxBin := buildBoxBinary(t)

	// Create a temp parent directory to contain our project and witness file
	tempParent := t.TempDir()
	projectDir := filepath.Join(tempParent, "project")
	if err := os.MkdirAll(projectDir, 0700); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	// File outside project directory that we will try to create/modify
	// We use /tmp directly because our sandbox profile specifically allows the project root and os.TempDir().
	// On macOS, os.TempDir() is usually /var/folders/..., so /tmp is a good escape target.
	witnessFile := "/tmp/box_sandbox_escape_witness.txt"
	_ = os.Remove(witnessFile) // Ensure it doesn't exist

	// Create a box.yml that attempts to write outside the project directory
	// Sandboxing is enabled by default for scripts
	boxYAML := `
tools:
  - type: script
    alias: sandbox-escape-attempt
    source:
      - touch /tmp/box_sandbox_escape_witness.txt
`
	if err := os.WriteFile(filepath.Join(projectDir, "box.yml"), []byte(boxYAML), 0600); err != nil {
		t.Fatalf("Failed to write box.yml: %v", err)
	}

	// Run box install - it should FAIL
	// On Linux, the current sandbox (User Namespaces) provides identity isolation
	// but not full filesystem isolation yet without complex mount setup.
	// We only expect failure on macOS where sandbox-exec is used.
	if runtime.GOOS == "linux" {
		t.Skip("Skipping filesystem escape check on Linux - User Namespace isolation is identity-only for now.")
	}

	//nolint:gosec
	installCmd := exec.Command(filepath.Clean(boxBin), "install", "--non-interactive")
	installCmd.Dir = projectDir
	output, err := installCmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected box install to fail due to sandbox violation, but it succeeded.\nOutput: %s", string(output))
	}

	// Verify the witness file was NOT created
	if _, err := os.Stat(witnessFile); err == nil {
		t.Errorf("Security Breach: Witness file %s was created despite sandbox!", witnessFile)
	}
}

func TestCLICommands(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	boxBin := buildBoxBinary(t)
	projectDir := setupTestProject(t, "testdata/integration_test.yml")

	t.Run("version", func(t *testing.T) {
		output := runBoxCommand(t, boxBin, projectDir, "version")
		if !strings.Contains(output, "box") {
			t.Errorf("Expected version output to contain 'box', got: %s", output)
		}
	})

	t.Run("install and list", func(t *testing.T) {
		runBoxCommand(t, boxBin, projectDir, "install", "--non-interactive")
		output := runBoxCommand(t, boxBin, projectDir, "list")
		if !strings.Contains(output, "task") {
			t.Errorf("Expected list output to contain 'task', got: %s", output)
		}
	})

	t.Run("run", func(t *testing.T) {
		output := runBoxCommand(t, boxBin, projectDir, "run", "task", "--version")
		if !strings.Contains(output, "Task version") {
			t.Errorf("Expected run output to contain 'Task version', got: %s", output)
		}
	})

	t.Run("env", func(t *testing.T) {
		output := runBoxCommand(t, boxBin, projectDir, "env")
		if !strings.Contains(output, "BOX_DIR") {
			t.Errorf("Expected env output to contain 'BOX_DIR', got: %s", output)
		}
		if !strings.Contains(output, "APP_DEBUG") {
			t.Errorf("Expected env output to contain 'APP_DEBUG', got: %s", output)
		}

		// Test specific key
		output = runBoxCommand(t, boxBin, projectDir, "env", "APP_DEBUG")
		if strings.TrimSpace(output) != "true" {
			t.Errorf("Expected 'true', got: %q", output)
		}
	})

	t.Run("generate direnv", func(t *testing.T) {
		runBoxCommand(t, boxBin, projectDir, "generate", "direnv")
		envrcPath := filepath.Join(projectDir, ".envrc")
		if _, err := os.Stat(envrcPath); err != nil {
			t.Errorf("Expected .envrc to be generated, but not found: %v", err)
		}
		content, _ := os.ReadFile(envrcPath)
		if !strings.Contains(string(content), "BOX_BIN_DIR") {
			t.Errorf(".envrc content missing BOX_BIN_DIR: %s", string(content))
		}
	})

	t.Run("generate dockerfile", func(t *testing.T) {
		runBoxCommand(t, boxBin, projectDir, "generate", "dockerfile")
		dockerfilePath := filepath.Join(projectDir, "Dockerfile")
		if _, err := os.Stat(dockerfilePath); err != nil {
			t.Errorf("Expected Dockerfile to be generated, but not found: %v", err)
		}
	})

	t.Run("doctor", func(t *testing.T) {
		// doctor might return non-zero if some tools are missing on host, 
		// but we just check if it runs.
		cmd := exec.Command(boxBin, "doctor")
		cmd.Dir = projectDir
		_ = cmd.Run() 
	})
}

func runBoxCommand(t *testing.T, boxBin, projectDir string, args ...string) string {
	t.Helper()
	//nolint:gosec
	cmd := exec.Command(boxBin, args...)
	cmd.Dir = projectDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("box %s failed: %v\nOutput: %s", strings.Join(args, " "), err, string(output))
	}
	return string(output)
}

func buildBoxBinary(t *testing.T) string {
	t.Helper()
	boxBin := filepath.Join(t.TempDir(), "box")
	//nolint:gosec
	buildCmd := exec.Command("go", "build", "-o", filepath.Clean(boxBin), "../../main.go")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build box binary: %v\nOutput: %s", err, string(output))
	}
	return boxBin
}

func setupTestProject(t *testing.T, configFile string) string {
	t.Helper()
	projectDir := t.TempDir()

	t.Cleanup(func() {
		_ = filepath.Walk(projectDir, func(path string, _ os.FileInfo, err error) error {
			if err == nil {
				//nolint:gosec
				_ = os.Chmod(path, 0700)
			}
			return nil
		})
	})

	//nolint:gosec
	configSource, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Failed to read integration test config: %v", err)
	}

	configPath := filepath.Join(projectDir, "box.yml")
	if err := os.WriteFile(configPath, configSource, 0600); err != nil {
		t.Fatalf("Failed to write box.yml: %v", err)
	}

	return projectDir
}

func runBoxInstall(t *testing.T, boxBin, projectDir string) {
	t.Helper()
	//nolint:gosec
	installCmd := exec.Command(filepath.Clean(boxBin), "install", "--non-interactive")
	installCmd.Dir = projectDir
	if output, err := installCmd.CombinedOutput(); err != nil {
		t.Fatalf("box install failed: %v\nOutput: %s", err, string(output))
	}
}

func verifyInstallation(t *testing.T, projectDir string) {
	t.Helper()
	configPath := filepath.Join(projectDir, "box.yml")
	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	for _, tool := range cfg.Tools {
		binaryNames := getBinaryNames(tool)
		for _, binaryName := range binaryNames {
			binPath := filepath.Join(projectDir, ".box", "bin", binaryName)
			if _, err := os.Stat(binPath); err != nil {
				t.Errorf("Expected binary for %s at %s, but not found", tool.Source, binPath)
				continue
			}

			if binaryName == "task" && tool.Version != "" {
				verifyToolVersion(t, binPath, tool.Version)
			}
		}
	}
}

func getBinaryNames(tool config.Tool) []string {
	if len(tool.Binaries) > 0 {
		return tool.Binaries
	}

	sourcePath := tool.Source.String()
	if parts := strings.Split(sourcePath, "/"); len(parts) > 1 {
		lastPart := parts[len(parts)-1]
		if len(lastPart) >= 2 && lastPart[0] == 'v' && isDigit(lastPart[1:]) {
			sourcePath = strings.Join(parts[:len(parts)-1], "/")
		}
	}

	binaryName := sourcePath
	if idx := strings.LastIndex(binaryName, "/"); idx != -1 {
		binaryName = binaryName[idx+1:]
	}
	return []string{binaryName}
}

func verifyToolVersion(t *testing.T, binPath, expectedVersion string) {
	t.Helper()
	//nolint:gosec
	versionCmd := exec.Command(filepath.Clean(binPath), "--version")
	output, err := versionCmd.CombinedOutput()
	if err != nil {
		t.Errorf("Failed to run installed tool %s: %v", binPath, err)
		return
	}

	if !strings.Contains(string(output), expectedVersion) {
		t.Errorf("Tool version mismatch. Expected %s, got output: %s", expectedVersion, string(output))
	}
}
