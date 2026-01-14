package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create a temporary etc.yml
	content := []byte(`
tools:
  - name: test-tool
    type: go
    source: example.com/tool
`)
	tmpfile, err := os.CreateTemp("", "etc-test-*.yml")
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

	if len(cfg.Tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(cfg.Tools))
	}

	if cfg.Tools[0].Name != "test-tool" {
		t.Errorf("Expected name 'test-tool', got '%s'", cfg.Tools[0].Name)
	}
}
