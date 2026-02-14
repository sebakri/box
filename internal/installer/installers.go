package installer

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sebakri/box/internal/config"
)

// Installer is the interface that all tool installers must implement.
type Installer interface {
	// Install installs the tool.
	Install(tool config.Tool, m *Manager) error
}

// runCommand is a helper to run shell commands with consistent output redirection and environment setup.
func (m *Manager) runCommand(name string, args []string, env []string, dir string) error {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdout = m.Output
	cmd.Stderr = m.Output

	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}

	m.log("Running: %s %s", name, strings.Join(args, " "))
	return cmd.Run()
}

// GoInstaller implements the Installer interface for Go tools.
type GoInstaller struct{}

// Install installs a Go tool using 'go install'.
func (i *GoInstaller) Install(tool config.Tool, m *Manager) error {
	boxDir := filepath.Join(m.RootDir, ".box")
	binDir := filepath.Join(boxDir, "bin")
	return m.installGo(tool, binDir)
}

// NpmInstaller implements the Installer interface for NPM packages.
type NpmInstaller struct{}

// Install installs an NPM package.
func (i *NpmInstaller) Install(tool config.Tool, m *Manager) error {
	boxDir := filepath.Join(m.RootDir, ".box")
	return m.installNpm(tool, boxDir)
}

// CargoInstaller implements the Installer interface for Cargo crates.
type CargoInstaller struct{}

// Install installs a Cargo crate using 'cargo-binstall'.
func (i *CargoInstaller) Install(tool config.Tool, m *Manager) error {
	boxDir := filepath.Join(m.RootDir, ".box")
	return m.installCargo(tool, boxDir)
}

// UvInstaller implements the Installer interface for Python tools via uv.
type UvInstaller struct{}

// Install installs a Python tool using 'uv tool install'.
func (i *UvInstaller) Install(tool config.Tool, m *Manager) error {
	boxDir := filepath.Join(m.RootDir, ".box")
	binDir := filepath.Join(boxDir, "bin")
	return m.installUv(tool, binDir)
}

// GemInstaller implements the Installer interface for Ruby gems.
type GemInstaller struct{}

// Install installs a Ruby gem.
func (i *GemInstaller) Install(tool config.Tool, m *Manager) error {
	boxDir := filepath.Join(m.RootDir, ".box")
	binDir := filepath.Join(boxDir, "bin")
	return m.installGem(tool, binDir)
}

// ScriptInstaller implements the Installer interface for shell scripts.
type ScriptInstaller struct{}

// Install installs a tool by running a shell script.
func (i *ScriptInstaller) Install(tool config.Tool, m *Manager) error {
	return m.installScript(tool)
}
