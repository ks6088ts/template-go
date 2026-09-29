package sandbox

import (
	"bytes"
	"testing"
)

// TestGetCommand checks the sandbox output and its nested command.
func TestGetCommand(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "sandbox", want: "sandbox called\n"},
		{name: "hello", args: []string{"hello"}, want: "sandbox hello called\n"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			command := GetCommand()
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
