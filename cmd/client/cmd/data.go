package cmd

import (
	"encoding/json"
	"github.com/Mldlr/storety/internal/client/models"
	"github.com/Mldlr/storety/internal/client/service"
	cobra "github.com/spf13/cobra"
	"io"
	"log"
	"os"
)

func dataClientCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "data",
		Short: "Data service operations",
		Long:  "Creating and deleting data",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	return cmd
}

func createCredentials(client *service.DataClient, crypto *service.Crypto) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create_cred [data_name] [login] [password] [meta]",
		Short: "Store new credentials pair",
		Long:  "",
		Args:  cobra.ExactArgs(4),
		RunE:  runCreateCredentials(client, crypto),
	}
	return cmd
}

func createCard(client *service.DataClient, crypto *service.Crypto) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create_card [data_name] [number] [expires] [name] [surname] [cvv] [meta]",
		Short: "Store new card",
		Long:  "",
		Args:  cobra.ExactArgs(7),
		RunE:  runCreateCard(client, crypto),
	}
	return cmd
}

func createText(client *service.DataClient, crypto *service.Crypto) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create_text [data_name] [text] [meta]",
		Short: "Store new text",
		Long:  "",
		Args:  cobra.ExactArgs(3),
		RunE:  runCreateText(client, crypto),
	}
	return cmd
}

func createBinary(client *service.DataClient, crypto *service.Crypto) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create_binary [data_name] [filename] [meta]",
		Short: "Store new binary",
		Long:  "",
		Args:  cobra.ExactArgs(3),
		RunE:  runCreateBinary(client, crypto),
	}
	return cmd
}

func listData(client *service.DataClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list_data",
		Short: "List data",
		Long:  "",
		Args:  cobra.ExactArgs(0),
		RunE:  runListData(client),
	}
	return cmd
}

func getData(client *service.DataClient, crypto *service.Crypto) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get_data [data_name]",
		Short: "Get data item",
		Long:  "",
		Args:  cobra.ExactArgs(1),
		RunE:  runGetData(client, crypto),
	}
	return cmd
}

func deleteData(client *service.DataClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete_data [data_name]",
		Short: "Get data item",
		Long:  "",
		Args:  cobra.ExactArgs(1),
		RunE:  runDeleteData(client),
	}
	return cmd
}

func runCreateCredentials(client *service.DataClient, crypto *service.Crypto) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		dataName := args[0]
		cred := &models.Credentials{
			Login:    args[1],
			Password: args[2],
			Meta:     args[3],
		}
		encodedCred, err := json.Marshal(cred)
		if err != nil {
			return logError(err)
		}
		encryptCred, err := crypto.EncryptWithAES256(encodedCred)
		err = client.CreateData(dataName, "Cred", encryptCred)
		if err != nil {
			return logError(err)
		}
		log.Println("Successfully created new credentials pair")
		return nil
	}
}

func runCreateCard(client *service.DataClient, crypto *service.Crypto) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
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
			return logError(err)
		}
		encryptCred, err := crypto.EncryptWithAES256(encodedCred)
		err = client.CreateData(dataName, "Card", encryptCred)
		if err != nil {
			return logError(err)
		}
		log.Println("Successfully created new card")
		return nil
	}
}

func runCreateText(client *service.DataClient, crypto *service.Crypto) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		dataName := args[0]
		cred := &models.Text{
			Text: args[1],
			Meta: args[2],
		}
		encodedCred, err := json.Marshal(cred)
		if err != nil {
			return logError(err)
		}
		encryptCred, err := crypto.EncryptWithAES256(encodedCred)
		err = client.CreateData(dataName, "Text", encryptCred)
		if err != nil {
			return logError(err)
		}
		log.Println("Successfully created new text")
		return nil
	}
}

func runCreateBinary(client *service.DataClient, crypto *service.Crypto) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		dataName := args[0]

		file, err := os.Open(args[1])
		if err != nil {
			return logError(err)
		}
		defer file.Close()

		blob, err := io.ReadAll(file)
		if err != nil {
			return logError(err)
		}
		cred := &models.Binary{
			Blob: blob,
			Meta: args[2],
		}
		encodedCred, err := json.Marshal(cred)
		if err != nil {
			return logError(err)
		}
		encryptCred, err := crypto.EncryptWithAES256(encodedCred)
		err = client.CreateData(dataName, "Binary", encryptCred)
		if err != nil {
			return logError(err)
		}
		log.Println("Successfully created new binary blob")
		return nil
	}
}

func runListData(client *service.DataClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		data, err := client.ListData()
		if err != nil {
			return logError(err)
		}
		for i, v := range data.Data {
			log.Printf("%d. %s - %s\n", i+1, v.Name, v.Type)
		}
		return nil
	}
}

func runGetData(client *service.DataClient, crypto *service.Crypto) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		data, err := client.GetData(args[0])
		if err != nil {
			return logError(err)
		}
		decryptCred, err := crypto.DecryptWithAES256(data.Content)
		if err != nil {
			return logError(err)
		}
		switch data.Type {
		case "Cred":
			cred := &models.Credentials{}
			err = json.Unmarshal(decryptCred, cred)
			if err != nil {
				return logError(err)
			}
			log.Printf("Login: %s\n", cred.Login)
			log.Printf("Password: %s\n", cred.Password)
			log.Printf("Meta: %s\n", cred.Meta)
		case "Card":
			card := &models.Card{}
			err = json.Unmarshal(decryptCred, card)
			if err != nil {
				return logError(err)
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
				return logError(err)
			}
			log.Printf("Text: %s\n", text.Text)
			log.Printf("Meta: %s\n", text.Meta)
		case "Binary":
			binary := &models.Binary{}
			err = json.Unmarshal(decryptCred, binary)
			if err != nil {
				return logError(err)
			}
			log.Printf("Blob: %s\n", binary.Blob)
			log.Printf("Blob written to file: %s\n", args[0])
			log.Printf("Meta: %s\n", binary.Meta)
			file, err := os.OpenFile(args[0], os.O_WRONLY|os.O_CREATE, 0755)
			if err != nil {
				return logError(err)
			}
			defer file.Close()
			_, err = file.Write(binary.Blob)
			if err != nil {
				return logError(err)
			}
		}
		return nil
	}
}

func runDeleteData(client *service.DataClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		err := client.DeleteData(args[0])
		if err != nil {
			return logError(err)
		}
		log.Println("Successfully deleted data")
		return nil
	}
}
