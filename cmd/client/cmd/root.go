package cmd

import (
	"context"
	"github.com/Mldlr/storety/internal/client/config"
	"github.com/Mldlr/storety/internal/client/service"
	shell "github.com/brianstrauch/cobra-shell"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type RunEFunc func(cmd *cobra.Command, args []string) error

var rootCmd = &cobra.Command{}

func Execute() {
	ctx := context.Background()
	cfg := config.NewConfig()
	conn, err := grpc.Dial(cfg.ServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logError(err)
	}
	defer conn.Close()
	crypto := service.NewCrypto(cfg)
	userClient := service.NewUserClient(ctx, conn, cfg)
	userCmd := userClientCommand()
	cUserCmd := createUserCmd(userClient)
	lUserCmd := logInCmd(userClient)
	userCmd.AddCommand(cUserCmd)
	userCmd.AddCommand(lUserCmd)
	rootCmd.AddCommand(userCmd)

	dataClient := service.NewDataClient(ctx, conn, cfg)
	dataCmd := dataClientCommand()
	dataCmd.AddCommand(createCredentials(dataClient, crypto))
	dataCmd.AddCommand(createCard(dataClient, crypto))
	dataCmd.AddCommand(createText(dataClient, crypto))
	dataCmd.AddCommand(createBinary(dataClient, crypto))
	dataCmd.AddCommand(listData(dataClient))
	dataCmd.AddCommand(getData(dataClient, crypto))
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
