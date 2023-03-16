// Package cmd contains the commands for the Storety client's command-line interface.
package cmd

import (
	"github.com/Mldlr/storety/internal/client/service"
	"github.com/Mldlr/storety/internal/client/service/crypto"
	shell "github.com/brianstrauch/cobra-shell"
	"github.com/spf13/cobra"
	"log"
)

// RunEFunc is a function type that can be used as a cobra command's RunE function.
type RunEFunc func(cmd *cobra.Command, args []string) error

// rootCmd is the root command of the Storety client's command-line interface.
var rootCmd = &cobra.Command{}

// Execute initializes and runs the root command along with its subcommands.
func Execute(userClient *service.UserClient, dataClient *service.DataClient, crypto *crypto.Crypto) {
	userCmd := userClientCommand()
	userCmd.AddCommand(logInCmd(userClient))
	userCmd.AddCommand(createUserCmd(userClient))
	rootCmd.AddCommand(userCmd)
	dataCmd := dataClientCommand()
	dataCmd.AddCommand(createCredentials(dataClient, crypto))
	dataCmd.AddCommand(createCard(dataClient, crypto))
	dataCmd.AddCommand(createText(dataClient, crypto))
	dataCmd.AddCommand(createBinary(dataClient, crypto))
	dataCmd.AddCommand(listData(dataClient))
	dataCmd.AddCommand(getData(dataClient, crypto))
	dataCmd.AddCommand(deleteData(dataClient))
	rootCmd.AddCommand(dataCmd)
	rootCmd.AddCommand(shell.New(rootCmd, nil))
	_ = rootCmd.Execute()
}

// logError logs an error if it is not nil.
func logError(err error) error {
	if err != nil {
		log.Println("Error running command: ", err)
	}
	return err
}
