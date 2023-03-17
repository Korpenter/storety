package cmd

import (
	"fmt"
	"github.com/Mldlr/storety/internal/client/service"
	cobra "github.com/spf13/cobra"
	"log"
)

// userClientCommand creates a cobra command for interacting with the user service.
func userClientCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "User service operations",
		Long:  "Registration and authentication",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	return cmd
}

// createUserCmd creates a cobra command for creating a new user.
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

// logInCmd creates a cobra command for logging in a user.
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

// runCreateUserCmd returns a RunEFunc that serves as a CLI wrapper for client.CreateUser.
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

// runLogInCmd returns a RunEFunc that serves as a CLI wrapper for client.LogInUser.
func runLogInCmd(client *service.UserClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		username, password := args[0], args[1]
		err := client.LogInUser(username, password)
		if err != nil {
			log.Println("Failed to log in on remote server, attempting local login")
			err = client.LocalLogin(username, password)
			if err != nil {
				return fmt.Errorf("failed local login: %v", err)
			}
			log.Println("Successful local log in")
			return nil
		}
		log.Println("Successful remote log in")
		return nil
	}
}
