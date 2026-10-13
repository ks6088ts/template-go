package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestCLI checks the real executable's streams and exit status without SetOut.
func TestCLI(t *testing.T) {
	executableName := "template-go"
	if runtime.GOOS == "windows" {
		executableName += ".exe"
	}
	executablePath := filepath.Join(t.TempDir(), executableName)
	buildCommand := exec.Command("go", "build", "-mod=readonly", "-o", executablePath, ".")
	if output, err := buildCommand.CombinedOutput(); err != nil {
		t.Fatalf("build CLI: %v\n%s", err, output)
	}
	tests := []struct {
		name     string
		args     []string
		want     string
		wantExit int
	}{
		{name: "version", args: []string{"version"}, want: "version=dev, revision=unknown\n"},
		{name: "sandbox", args: []string{"sandbox"}, want: "sandbox called\n"},
		{name: "hello", args: []string{"sandbox", "hello"}, want: "sandbox hello called\n"},
		{name: "invalid argument", args: []string{"version", "extra"}, wantExit: 1},
		{name: "unknown command", args: []string{"missing"}, wantExit: 1},
		{name: "removed flag", args: []string{"--config", "missing.yaml"}, wantExit: 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			var diagnostic bytes.Buffer
			command := exec.Command(executablePath, test.args...)
			command.Stdout = &output
			command.Stderr = &diagnostic
			err := command.Run()
			if test.wantExit == 0 {
				if err != nil {
					t.Fatalf("run CLI: %v\n%s", err, &diagnostic)
				}
				if diagnostic.Len() != 0 {
					t.Errorf("unexpected stderr: %q", diagnostic.String())
				}
			} else {
				var exitError *exec.ExitError
				if !errors.As(err, &exitError) || exitError.ExitCode() != test.wantExit {
					t.Fatalf("run error = %v, want exit %d", err, test.wantExit)
				}
				if strings.Count(diagnostic.String(), "Error:") != 1 || strings.Contains(diagnostic.String(), "Usage:") {
					t.Errorf("stderr = %q, want one error without usage", diagnostic.String())
				}
			}
			if output.String() != test.want {
				t.Errorf("stdout = %q, want %q", output.String(), test.want)
			}
		})
	}
	if runtime.GOOS == "linux" {
		t.Run("write failure", func(t *testing.T) {
			// /dev/full fails writes deterministically, unlike a pipe closed by a reader.
			output, err := os.OpenFile("/dev/full", os.O_WRONLY, 0)
			if err != nil {
				t.Fatalf("open failing output: %v", err)
			}
			t.Cleanup(func() {
				if err := output.Close(); err != nil {
					t.Errorf("close output: %v", err)
				}
			})
			var diagnostic bytes.Buffer
			command := exec.Command(executablePath, "version")
			command.Stdout = output
			command.Stderr = &diagnostic
			var exitError *exec.ExitError
			if err := command.Run(); !errors.As(err, &exitError) || exitError.ExitCode() != 1 {
				t.Fatalf("run error = %v, want exit 1", err)
			}
			if !strings.Contains(diagnostic.String(), "Error: write version:") {
				t.Errorf("stderr = %q, want write error", diagnostic.String())
			}
		})
	}
}

// TestReleaseTag checks the shared release gate without publishing anything.
func TestReleaseTag(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Make release checks run on Unix; Windows only needs Go for development")
	}
	makePath, err := exec.LookPath("make")
	if err != nil {
		t.Skip("release tooling tests require GNU Make")
	}
	tests := []struct {
		tag   string
		valid bool
	}{
		{tag: "v0.0.0", valid: true},
		{tag: "v1.2.3", valid: true},
		{tag: "v1.2.3-rc.1", valid: true},
		{tag: "v1.2.3-dirty", valid: true},
		{tag: "v1.2.3-0", valid: true},
		{tag: "v1.2.3-01alpha", valid: true},
		{tag: "v1.2.3+build.01", valid: true},
		{tag: "v1.2.3-rc.1+build.5", valid: true},
		{tag: ""},
		{tag: "1.2.3"},
		{tag: "v1.2"},
		{tag: "v01.2.3"},
		{tag: "v1.02.3"},
		{tag: "v1.2.03"},
		{tag: "v1.2.3-01"},
		{tag: "v1.2.3-rc.01"},
		{tag: "v1.2.3-"},
		{tag: "v1.2.3+"},
		{tag: "v1.2.3-rc..1"},
		{tag: "v1.2.3+build..1"},
		{tag: "v1.2.3 extra"},
		{tag: "v1.2.3\nextra"},
		{tag: "v1.2.3\r"},
	}
	for _, test := range tests {
		t.Run(test.tag, func(t *testing.T) {
			command := exec.Command(makePath, "--no-print-directory", "release-check")
			command.Env = append(os.Environ(), "RELEASE_TAG="+test.tag)
			output, err := command.CombinedOutput()
			if test.valid {
				if err != nil {
					t.Fatalf("valid tag rejected: %v\n%s", err, output)
				}
			} else {
				var exitError *exec.ExitError
				if !errors.As(err, &exitError) || !bytes.Contains(output, []byte("Invalid release tag:")) {
					t.Fatalf("invalid tag: error = %v, output = %q", err, output)
				}
			}
		})
	}
}
