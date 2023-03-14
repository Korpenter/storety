package cmd

import (
	"github.com/Mldlr/storety/internal/client/service"
	cobra "github.com/spf13/cobra"
	"log"
)

func userClientCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "User service operations",
		Long:  "Registration and authentication",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	return cmd
}

func createUserCmd(client *service.UserClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [name] [password]",
		Short: "Create new account",
		Long:  "",
		Args:  cobra.ExactArgs(2),
		RunE:  runCreateUserCmd(client),
	}
	return cmd
}

func logInCmd(client *service.UserClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [name] [password]",
		Short: "Log into account",
		Long:  "",
		Args:  cobra.ExactArgs(2),
		RunE:  runLogInCmd(client),
	}
	return cmd
}

func runCreateUserCmd(client *service.UserClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		username, password := args[0], args[1]
		err := client.CreateUser(username, password)
		if err != nil {
			return logError(err)
		}
		log.Println("Successfully created new user")
		return nil
	}
}

func runLogInCmd(client *service.UserClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		username, password := args[0], args[1]
		err := client.LogInUser(username, password)
		if err != nil {
			return logError(err)
		}
		log.Println("Successful log in")
		return nil
	}
}
