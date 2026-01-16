package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create a temporary box.yml
	content := []byte(`
tools:
  - type: go
    source: example.com/tool
  - type: uv
    source: ruff
env:
  KEY: value
`)
	tmpfile, err := os.CreateTemp("", "box-test-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test loading
	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(cfg.Tools) != 2 {
		t.Errorf("Expected 2 tools, got %d", len(cfg.Tools))
	}

	if cfg.Env["KEY"] != "value" {
		t.Errorf("Expected Env['KEY'] to be 'value', got '%s'", cfg.Env["KEY"])
	}

	if cfg.Tools[0].Source != "example.com/tool" {
		t.Errorf("Expected source 'example.com/tool', got '%s'", cfg.Tools[0].Source)
	}
}
