package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ks6088ts/template-go/internal"
)

// TestCommandOutput checks the user-facing commands without process globals.
func TestCommandOutput(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "version", args: []string{"version"}, want: fmt.Sprintf("version=%s, revision=%s\n", internal.Version, internal.Revision)},
		{name: "sandbox", args: []string{"sandbox"}, want: "sandbox called\n"},
		{name: "nested command", args: []string{"sandbox", "hello"}, want: "sandbox hello called\n"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			command := newRootCommand()
			command.SetArgs(test.args)
			command.SetOut(&output)
			if err := command.Execute(); err != nil {
				t.Fatalf("execute %v: %v", test.args, err)
			}
			if got := output.String(); got != test.want {
				t.Errorf("output = %q, want %q", got, test.want)
			}
		})
	}
}

// TestConfig verifies that only an absent default config is optional.
func TestConfig(t *testing.T) {
	tests := []struct {
		name         string
		contents     string
		createFile   bool
		explicitFile bool
		wantError    bool
	}{
		{name: "absent default"},
		{name: "valid default", contents: "key: value\n", createFile: true},
		{name: "invalid default", contents: "key: [\n", createFile: true, wantError: true},
		{name: "valid explicit", contents: "key: value\n", createFile: true, explicitFile: true},
		{name: "missing explicit", explicitFile: true, wantError: true},
		{name: "invalid explicit", contents: "key: [\n", createFile: true, explicitFile: true, wantError: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			configFile := filepath.Join(home, ".template-go.yaml")
			if test.createFile {
				if err := os.WriteFile(configFile, []byte(test.contents), 0600); err != nil {
					t.Fatalf("write config: %v", err)
				}
			}

			args := []string{"version"}
			if test.explicitFile {
				args = append(args, "--config", configFile)
			}
			var output bytes.Buffer
			var errors bytes.Buffer
			command := newRootCommand()
			command.SetArgs(args)
			command.SetOut(&output)
			command.SetErr(&errors)
			err := command.Execute()
			if test.wantError {
				if err == nil || !strings.Contains(err.Error(), "read config:") {
					t.Fatalf("execute error = %v, want a config error", err)
				}
				if strings.Contains(output.String(), "version=") {
					t.Errorf("command ran on config failure: %q", output.String())
				}
				if !strings.Contains(errors.String(), "Error: read config:") {
					t.Errorf("diagnostic = %q, want config error", errors.String())
				}
				return
			}
			if err != nil {
				t.Fatalf("execute: %v", err)
			}
			if !strings.HasPrefix(output.String(), "version=") {
				t.Errorf("output = %q, want version", output.String())
			}
			if test.createFile && !strings.Contains(errors.String(), "Using config file: "+configFile) {
				t.Errorf("diagnostic = %q, want config path", errors.String())
			}
			if !test.createFile && errors.Len() != 0 {
				t.Errorf("unexpected diagnostic: %q", errors.String())
			}
		})
	}
}
