package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// TestBuildMetadata checks source-archive defaults and explicit build overrides.
func TestBuildMetadata(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Make builds are tested on Unix; Windows development only requires Go")
	}
	makePath, err := exec.LookPath("make")
	if err != nil {
		t.Skip("build tooling tests require GNU Make")
	}
	sourceDir := t.TempDir()
	for _, directory := range []string{"cmd", "internal"} {
		if err := os.CopyFS(filepath.Join(sourceDir, directory), os.DirFS(directory)); err != nil {
			t.Fatalf("copy %s: %v", directory, err)
		}
	}
	for _, fileName := range []string{"main.go", "go.mod", "go.sum", "Makefile"} {
		contents, err := os.ReadFile(fileName)
		if err != nil {
			t.Fatalf("read %s: %v", fileName, err)
		}
		if err := os.WriteFile(filepath.Join(sourceDir, fileName), contents, 0600); err != nil {
			t.Fatalf("copy %s: %v", fileName, err)
		}
	}
	tests := []struct {
		name      string
		overrides []string
		want      string
	}{
		{name: "source archive", want: "version=dev, revision=unknown\n"},
		{name: "explicit metadata", overrides: []string{"GIT_TAG=1.2.3", "GIT_REVISION=abc123"}, want: "version=1.2.3, revision=abc123\n"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			executablePath := filepath.Join(sourceDir, "app")
			args := []string{"build", "OUTPUT=" + executablePath, "GOOS=" + runtime.GOOS, "GOARCH=" + runtime.GOARCH}
			args = append(args, test.overrides...)
			buildCommand := exec.Command(makePath, args...)
			buildCommand.Dir = sourceDir
			if output, err := buildCommand.CombinedOutput(); err != nil {
				t.Fatalf("build: %v\n%s", err, output)
			}
			output, err := exec.Command(executablePath, "version").CombinedOutput()
			if err != nil {
				t.Fatalf("run CLI: %v\n%s", err, output)
			}
			if string(output) != test.want {
				t.Errorf("version = %q, want %q", output, test.want)
			}
		})
	}
}
