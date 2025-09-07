package container

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultEntrypointCmdFallback(t *testing.T) {
	os.Unsetenv("REPOCATE_ENTRYPOINT")
	cmd := defaultEntrypointCmd()
	if len(cmd) != 1 || cmd[0] != "/bin/zsh" {
		t.Fatalf("expected fallback to /bin/zsh, got %v", cmd)
	}
}

func TestDefaultEntrypointCmdEnv(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "entrypoint.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("failed to write temp script: %v", err)
	}
	os.Setenv("REPOCATE_ENTRYPOINT", script)
	defer os.Unsetenv("REPOCATE_ENTRYPOINT")

	cmd := defaultEntrypointCmd()
	if len(cmd) != 1 || cmd[0] != script {
		t.Fatalf("expected %s, got %v", script, cmd)
	}
}
