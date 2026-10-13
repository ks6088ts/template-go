package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ks6088ts/template-go/internal"
)

// TestCommandOutput checks the user-facing commands without process globals.
func TestCommandOutput(t *testing.T) {
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
			var diagnostic bytes.Buffer
			command := newRootCommand()
			command.SetArgs(test.args)
			command.SetOut(&output)
			command.SetErr(&diagnostic)
			if err := command.Execute(); err != nil {
				t.Fatalf("execute %v: %v", test.args, err)
			}
			if got := output.String(); got != test.want {
				t.Errorf("output = %q, want %q", got, test.want)
			}
			if diagnostic.Len() != 0 {
				t.Errorf("unexpected diagnostic: %q", diagnostic.String())
			}
		})
	}
}

// TestInvalidArguments verifies that invalid input never runs a command.
func TestInvalidArguments(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "unknown command", args: []string{"missing"}, want: "unknown command"},
		{name: "unknown flag", args: []string{"version", "--missing"}, want: "unknown flag"},
		{name: "removed config", args: []string{"version", "--config", "config.yaml"}, want: "unknown flag"},
		{name: "removed toggle", args: []string{"--toggle"}, want: "unknown flag"},
		{name: "version argument", args: []string{"version", "extra"}, want: "unknown command"},
		{name: "sandbox argument", args: []string{"sandbox", "extra"}, want: "unknown command"},
		{name: "hello argument", args: []string{"sandbox", "hello", "extra"}, want: "unknown command"},
		{name: "argument after separator", args: []string{"version", "--", "extra"}, want: "unknown command"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			var diagnostic bytes.Buffer
			command := newRootCommand()
			command.SetArgs(test.args)
			command.SetOut(&output)
			command.SetErr(&diagnostic)
			err := command.Execute()
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("execute error = %v, want %q", err, test.want)
			}
			if output.Len() != 0 {
				t.Errorf("unexpected output: %q", output.String())
			}
			if strings.Count(diagnostic.String(), "Error:") != 1 {
				t.Errorf("diagnostic = %q, want one error", diagnostic.String())
			}
			if strings.Contains(diagnostic.String(), "Usage:") {
				t.Errorf("unexpected usage: %q", diagnostic.String())
			}
		})
	}
}

// TestHelpAndCompletion checks Cobra's built-in commands after simplification.
func TestHelpAndCompletion(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "root", args: []string{}, want: "Usage:"},
		{name: "help flag", args: []string{"--help"}, want: "Usage:"},
		{name: "nested help", args: []string{"help", "sandbox", "hello"}, want: "template-go sandbox hello"},
		{name: "completion", args: []string{"completion", "bash"}, want: "__start_template-go"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			var diagnostic bytes.Buffer
			command := newRootCommand()
			command.SetArgs(test.args)
			command.SetOut(&output)
			command.SetErr(&diagnostic)
			if err := command.Execute(); err != nil {
				t.Fatalf("execute: %v", err)
			}
			if !strings.Contains(output.String(), test.want) {
				t.Errorf("output = %q, want %q", output.String(), test.want)
			}
			if diagnostic.Len() != 0 {
				t.Errorf("unexpected diagnostic: %q", diagnostic.String())
			}
		})
	}
}

// failingWriter simulates an output stream that cannot accept data.
type failingWriter struct {
	err error
}

// Write returns the configured failure without consuming data.
func (writer failingWriter) Write(data []byte) (int, error) {
	return 0, writer.err
}

// TestOutputFailure ensures every application command preserves writer errors.
func TestOutputFailure(t *testing.T) {
	writeError := errors.New("output unavailable")
	for _, args := range [][]string{{"version"}, {"sandbox"}, {"sandbox", "hello"}} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			var diagnostic bytes.Buffer
			command := newRootCommand()
			command.SetArgs(args)
			command.SetOut(failingWriter{err: writeError})
			command.SetErr(&diagnostic)
			if err := command.Execute(); !errors.Is(err, writeError) {
				t.Fatalf("execute error = %v, want %v", err, writeError)
			}
			if strings.Count(diagnostic.String(), "Error:") != 1 || strings.Contains(diagnostic.String(), "Usage:") {
				t.Errorf("diagnostic = %q, want one error without usage", diagnostic.String())
			}
		})
	}
}

// TestFreshCommandState checks that a failed invocation does not affect the next.
func TestFreshCommandState(t *testing.T) {
	var diagnostic bytes.Buffer
	firstCommand := newRootCommand()
	firstCommand.SetArgs([]string{"version", "--missing"})
	firstCommand.SetErr(&diagnostic)
	if err := firstCommand.Execute(); err == nil {
		t.Fatal("invalid flag unexpectedly succeeded")
	}
	var output bytes.Buffer
	secondCommand := newRootCommand()
	secondCommand.SetArgs([]string{"sandbox", "hello"})
	secondCommand.SetOut(&output)
	if err := secondCommand.Execute(); err != nil {
		t.Fatalf("fresh command: %v", err)
	}
	if output.String() != "sandbox hello called\n" {
		t.Errorf("output = %q", output.String())
	}
}
