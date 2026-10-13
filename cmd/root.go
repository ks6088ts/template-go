/*
Copyright © 2024 ks6088ts

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/ks6088ts/template-go/cmd/sandbox"
	"github.com/spf13/cobra"
)

// Execute runs a fresh command tree and returns any command error.
func Execute() error {
	return newRootCommand().Execute()
}

// newRootCommand isolates command state so each invocation starts clean.
func newRootCommand() *cobra.Command {
	rootCommand := &cobra.Command{
		Use:          "template-go",
		Short:        "A minimal Go CLI template",
		Example:      "  template-go sandbox hello\n  template-go version",
		SilenceUsage: true,
	}
	rootCommand.AddCommand(sandbox.GetCommand(), newVersionCommand())
	return rootCommand
}
