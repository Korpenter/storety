package cmd

import (
	"fmt"
	"github.com/Mldlr/storety/internal/client/config"
	"github.com/Mldlr/storety/internal/client/pkg/helpers"
	"github.com/Mldlr/storety/internal/client/service/data"
	"github.com/Mldlr/storety/internal/client/service/user"
	"github.com/Mldlr/storety/internal/client/storage/sqlite"
	"github.com/samber/do"
	cobra "github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
)

// userClientCommand creates a cobra command for interacting with the user service.
func userClientCommand(i *do.Injector) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "User service operations",
		Long:  "Registration and authentication",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	return cmd
}

// createUserCmd creates a cobra command for creating a new user.
func createUserCmd(i *do.Injector) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [name] [password]",
		Short: "Create new account",
		Long:  "",
		Args:  cobra.ExactArgs(2),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			conn := do.MustInvoke[*grpc.ClientConn](i)
			if conn == nil {
				return helpers.LogError(fmt.Errorf("NO SERVER CONNECTION"))
			}
			return nil
		},
		RunE: runCreateUserCmd(i),
	}
	return cmd
}

// logInCmd creates a cobra command for logging in a user.
func logInCmd(i *do.Injector) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [name] [password]",
		Short: "Log into account",
		Long:  "",
		Args:  cobra.ExactArgs(2),
		RunE:  runLogInCmd(i),
	}
	return cmd
}

// runCreateUserCmd returns a RunEFunc that serves as a CLI wrapper for client.CreateUser.
func runCreateUserCmd(i *do.Injector) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		cfg := do.MustInvoke[*config.Config](i)
		userService := do.MustInvoke[user.Service](i)
		dataService := do.MustInvoke[data.Service](i)
		username, password := args[0], args[1]
		err := userService.CreateUser(username, password)
		if err != nil {
			return helpers.LogError(err)
		}
		log.Println("Successfully created new user")
		db, err := sqlite.NewDB(cfg.DBFilePrefix, username)
		if err != nil {
			return helpers.LogError(fmt.Errorf("failed to create new local database, u can continue using remote: %v", err))
		}
		dataService.SetStorage(db)
		log.Println("Successfully initiated new local database")
		return nil
	}
}

// runLogInCmd returns a RunEFunc that serves as a CLI wrapper for client.LogInUser.
func runLogInCmd(i *do.Injector) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		userService := do.MustInvoke[user.Service](i)
		dataService := do.MustInvoke[data.Service](i)
		cfg := do.MustInvoke[*config.Config](i)
		username, password := args[0], args[1]
		err := userService.LogInUser(username, password)
		if err != nil {
			return helpers.LogError(err)
		}
		log.Println("Successfully logged in")
		db, err := sqlite.NewDB(cfg.DBFilePrefix, username)
		if err != nil {
			return helpers.LogError(fmt.Errorf("failed to locate or create local database: %v", err))
		}
		dataService.SetStorage(db)
		return nil
	}
}
