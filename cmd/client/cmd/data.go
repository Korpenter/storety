package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/Mldlr/storety/internal/client/config"
	"github.com/Mldlr/storety/internal/client/models"
	"github.com/Mldlr/storety/internal/client/pkg/helpers"
	"github.com/Mldlr/storety/internal/client/service/crypto"
	"github.com/Mldlr/storety/internal/client/service/data"
	"github.com/samber/do"
	cobra "github.com/spf13/cobra"
	"io"
	"log"
	"os"
)

// dataClientCommand creates a cobra command for interacting with data service.
func dataClientCommand(i *do.Injector) *cobra.Command {
	cfg := do.MustInvoke[*config.Config](i)
	cmd := &cobra.Command{
		Use:   "data",
		Short: "Data service operations",
		Long:  "Creating and deleting data",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cfg.EncryptionKey == nil {
				return helpers.LogError(fmt.Errorf("NOT LOGGED IN"))
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	return cmd
}

// createCredentials creates a cobra command for creating a new credentials pair.
func createCredentials(i *do.Injector) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create_cred [data_name] [login] [password] [meta]",
		Short: "Store new credentials pair",
		Long:  "",
		Args:  cobra.ExactArgs(4),
		RunE:  runCreateCredentials(i),
	}
	return cmd
}

// createCard creates a cobra command for creating a new card.
func createCard(i *do.Injector) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create_card [data_name] [number] [expires] [name] [surname] [cvv] [meta]",
		Short: "Store new card",
		Long:  "",
		Args:  cobra.ExactArgs(7),
		RunE:  runCreateCard(i),
	}
	return cmd
}

// createText creates a cobra command for creating a new text data item.
func createText(i *do.Injector) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create_text [data_name] [text] [meta]",
		Short: "Store new text",
		Long:  "",
		Args:  cobra.ExactArgs(3),
		RunE:  runCreateText(i),
	}
	return cmd
}

// createBinary creates a cobra command for creating a new binary data item.
func createBinary(i *do.Injector) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create_binary [data_name] [filename] [meta]",
		Short: "Store new binary",
		Long:  "",
		Args:  cobra.ExactArgs(3),
		RunE:  runCreateBinary(i),
	}
	return cmd
}

// listData creates a cobra command for listing all data items.
func listData(i *do.Injector) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List data",
		Long:  "",
		Args:  cobra.ExactArgs(0),
		RunE:  runListData(i),
	}
	return cmd
}

// getData creates a cobra command for getting a data item.
func getData(i *do.Injector) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [data_name]",
		Short: "Get data item",
		Long:  "",
		Args:  cobra.ExactArgs(1),
		RunE:  runGetData(i),
	}
	return cmd
}

// deleteData creates a cobra command for deleting a data item.
func deleteData(i *do.Injector) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [data_name]",
		Short: "Get data item",
		Long:  "",
		Args:  cobra.ExactArgs(1),
		RunE:  runDeleteData(i),
	}
	return cmd
}

// syncData creates a cobra command for deleting a data item.
func syncData(i *do.Injector) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Manual Sync with remote",
		Long:  "",
		Args:  cobra.ExactArgs(0),
		RunE:  runSync(i),
	}
	return cmd
}

// runCreateCredentials is a wrapper creating a new Credentials data item.
func runCreateCredentials(i *do.Injector) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		dataService := do.MustInvoke[data.Service](i)
		cryptoService := do.MustInvoke[crypto.Crypto](i)
		dataName := args[0]
		cred := &models.Credentials{
			Login:    args[1],
			Password: args[2],
			Meta:     args[3],
		}
		encodedCred, err := json.Marshal(cred)
		if err != nil {
			return helpers.LogError(err)
		}
		encryptCred, err := cryptoService.EncryptWithAES256(encodedCred)
		err = dataService.CreateData(dataName, "Cred", encryptCred)
		if err != nil {
			return helpers.LogError(err)
		}
		log.Println("Successfully created new credentials pair")
		return nil
	}
}

// runCreateCard is a wrapper creating a new Card data item.
func runCreateCard(i *do.Injector) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		dataService := do.MustInvoke[data.Service](i)
		cryptoService := do.MustInvoke[crypto.Crypto](i)
		dataName := args[0]
		cred := &models.Card{
			Number:  args[1],
			Expires: args[2],
			CVV:     args[3],
			Name:    args[4],
			Surname: args[5],
			Meta:    args[6],
		}
		encodedCred, err := json.Marshal(cred)
		if err != nil {
			return helpers.LogError(err)
		}
		encryptCred, err := cryptoService.EncryptWithAES256(encodedCred)
		err = dataService.CreateData(dataName, "Card", encryptCred)
		if err != nil {
			return helpers.LogError(err)
		}
		log.Println("Successfully created new card")
		return nil
	}
}

// runCreateText is a wrapper creating a new text data item.
func runCreateText(i *do.Injector) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		dataService := do.MustInvoke[data.Service](i)
		cryptoService := do.MustInvoke[crypto.Crypto](i)
		dataName := args[0]
		cred := &models.Text{
			Text: args[1],
			Meta: args[2],
		}
		encodedCred, err := json.Marshal(cred)
		if err != nil {
			return helpers.LogError(err)
		}
		encryptCred, err := cryptoService.EncryptWithAES256(encodedCred)
		err = dataService.CreateData(dataName, "Text", encryptCred)
		if err != nil {
			return helpers.LogError(err)
		}
		log.Println("Successfully created new text")
		return nil
	}
}

// runCreateBinary is a wrapper creating a new binary data item.
func runCreateBinary(i *do.Injector) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		dataService := do.MustInvoke[data.Service](i)
		cryptoService := do.MustInvoke[crypto.Crypto](i)
		dataName := args[0]
		file, err := os.Open(args[1])
		if err != nil {
			return helpers.LogError(err)
		}
		defer file.Close()

		blob, err := io.ReadAll(file)
		if err != nil {
			return helpers.LogError(err)
		}
		cred := &models.Binary{
			Blob: blob,
			Meta: args[2],
		}
		encodedCred, err := json.Marshal(cred)
		if err != nil {
			return helpers.LogError(err)
		}
		encryptCred, err := cryptoService.EncryptWithAES256(encodedCred)
		err = dataService.CreateData(dataName, "Binary", encryptCred)
		if err != nil {
			return helpers.LogError(err)
		}
		log.Println("Successfully created new binary blob")
		return nil
	}
}

// runListData is a wrapper for getting data info from the server
func runListData(i *do.Injector) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		dataService := do.MustInvoke[data.Service](i)
		data, err := dataService.ListData()
		if err != nil {
			return helpers.LogError(err)
		}
		for i, v := range data {
			log.Printf("%d. %s - %s\n", i+1, v.Name, v.Type)
		}
		return nil
	}
}

// runGetData is a wrapper for getting data from the server and formatting it.
func runGetData(i *do.Injector) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		dataService := do.MustInvoke[data.Service](i)
		cryptoService := do.MustInvoke[crypto.Crypto](i)
		content, typ, err := dataService.GetData(args[0])
		if err != nil {
			return helpers.LogError(err)
		}
		decryptCred, err := cryptoService.DecryptWithAES256(content)
		if err != nil {
			return helpers.LogError(err)
		}
		switch typ {
		case "Cred":
			cred := &models.Credentials{}
			err = json.Unmarshal(decryptCred, cred)
			if err != nil {
				return helpers.LogError(err)
			}
			log.Printf("Login: %s\n", cred.Login)
			log.Printf("Password: %s\n", cred.Password)
			log.Printf("Meta: %s\n", cred.Meta)
		case "Card":
			card := &models.Card{}
			err = json.Unmarshal(decryptCred, card)
			if err != nil {
				return helpers.LogError(err)
			}
			log.Printf("Number: %s\n", card.Number)
			log.Printf("Expires: %s\n", card.Expires)
			log.Printf("CVV: %s\n", card.CVV)
			log.Printf("Name: %s\n", card.Name)
			log.Printf("Surname: %s\n", card.Surname)
			log.Printf("Meta: %s\n", card.Meta)
		case "Text":
			text := &models.Text{}
			err = json.Unmarshal(decryptCred, text)
			if err != nil {
				return helpers.LogError(err)
			}
			log.Printf("Text: %s\n", text.Text)
			log.Printf("Meta: %s\n", text.Meta)
		case "Binary":
			binary := &models.Binary{}
			err = json.Unmarshal(decryptCred, binary)
			if err != nil {
				return helpers.LogError(err)
			}
			log.Printf("Blob: %s\n", binary.Blob)
			log.Printf("Meta: %s\n", binary.Meta)
		}
		return nil
	}
}

// runDeleteData is a wrapper for deleting data.
func runDeleteData(i *do.Injector) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		dataService := do.MustInvoke[data.Service](i)
		err := dataService.DeleteData(args[0])
		if err != nil {
			return helpers.LogError(err)
		}
		log.Println("Successfully deleted data")
		return nil
	}
}

// runSync is a wrapper for syncing data.
func runSync(i *do.Injector) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		dataService := do.MustInvoke[data.Service](i)
		err := dataService.SyncData()
		if err != nil {
			log.Println("Error syncing data: ", err)
			return helpers.LogError(err)
		}
		log.Println("Successfully synced data")
		return nil
	}
}
