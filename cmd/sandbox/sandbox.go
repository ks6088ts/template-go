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
package sandbox

import (
	"fmt"

	"github.com/spf13/cobra"
)

// GetCommand creates a sandbox command with its nested hello command.
func GetCommand() *cobra.Command {
	sandboxCommand := &cobra.Command{
		Use:   "sandbox",
		Short: "Run the example command",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, args []string) error {
			if _, err := fmt.Fprintln(command.OutOrStdout(), "sandbox called"); err != nil {
				return fmt.Errorf("write sandbox output: %w", err)
			}
			return nil
		},
	}
	sandboxCommand.AddCommand(newHelloCommand())
	return sandboxCommand
}
