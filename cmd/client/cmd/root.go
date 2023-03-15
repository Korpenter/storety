package cmd

import (
	"github.com/Mldlr/storety/internal/client/service"
	shell "github.com/brianstrauch/cobra-shell"
	"github.com/spf13/cobra"
	"log"
)

type RunEFunc func(cmd *cobra.Command, args []string) error

var rootCmd = &cobra.Command{}

func Execute(userClient *service.UserClient, dataClient *service.DataClient, crypto *service.Crypto) {
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

func init() {
	rootCmd.PersistentFlags().StringP("server", "s", "localhost:8081", "set grpc server address")
}
func logError(err error) error {
	if err != nil {
		log.Println("Error running command: ", err)
	}
	return err
}
