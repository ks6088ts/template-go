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
	"errors"
	"fmt"
	"os"

	"github.com/ks6088ts/template-go/cmd/sandbox"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Execute runs a fresh command tree and returns any command error.
func Execute() error {
	return newRootCommand().Execute()
}

// newRootCommand isolates flags and configuration so each invocation starts clean.
func newRootCommand() *cobra.Command {
	var configFile string
	config := viper.New()
	rootCommand := &cobra.Command{
		Use:   "template-go",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}
	rootCommand.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.template-go.yaml)")
	rootCommand.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCommand.PersistentPreRunE = func(command *cobra.Command, args []string) error {
		if configFile != "" {
			config.SetConfigFile(configFile)
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("find home directory: %w", err)
			}
			config.AddConfigPath(home)
			config.SetConfigType("yaml")
			config.SetConfigName(".template-go")
		}
		config.AutomaticEnv()

		if err := config.ReadInConfig(); err != nil {
			var notFound viper.ConfigFileNotFoundError
			if configFile == "" && errors.As(err, &notFound) {
				return nil
			}
			return fmt.Errorf("read config: %w", err)
		}
		if _, err := fmt.Fprintln(command.ErrOrStderr(), "Using config file:", config.ConfigFileUsed()); err != nil {
			return fmt.Errorf("report config file: %w", err)
		}
		return nil
	}
	rootCommand.AddCommand(sandbox.GetCommand(), newVersionCommand())
	return rootCommand
}
