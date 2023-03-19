// Package cmd contains the commands for the Storety client's command-line interface.
package cmd

import (
	shell "github.com/brianstrauch/cobra-shell"
	"github.com/samber/do"
	"github.com/spf13/cobra"
)

// RunEFunc is a function type that can be used as a cobra command's RunE function.
type RunEFunc func(cmd *cobra.Command, args []string) error

// rootCmd is the root command of the Storety client's command-line interface.
var rootCmd = &cobra.Command{}

// Execute initializes and runs the root command along with its subcommands.
func Execute(i *do.Injector) {
	userCmd := userClientCommand(i)
	userCmd.AddCommand(logInCmd(i))
	userCmd.AddCommand(createUserCmd(i))
	rootCmd.AddCommand(userCmd)
	dataCmd := dataClientCommand(i)
	dataCmd.AddCommand(createCredentials(i))
	dataCmd.AddCommand(createCard(i))
	dataCmd.AddCommand(createText(i))
	dataCmd.AddCommand(createBinary(i))
	dataCmd.AddCommand(listData(i))
	dataCmd.AddCommand(getData(i))
	dataCmd.AddCommand(deleteData(i))
	dataCmd.AddCommand(syncData(i))
	rootCmd.AddCommand(dataCmd)
	rootCmd.AddCommand(shell.New(rootCmd, nil))
	_ = rootCmd.Execute()
}
