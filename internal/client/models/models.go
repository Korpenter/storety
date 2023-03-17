// Package models contains all the models used by the client.
package models

// Credentials is a struct that represents a credentials pair.
type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Meta     string `json:"meta"`
}

// Card is a struct that represents a card.
type Card struct {
	Number  string `json:"number"`
	Expires string `json:"expires"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	CVV     string `json:"cvv"`
	Meta    string `json:"meta"`
}

// Text is a struct that represents a text.
type Text struct {
	Text string `json:"text"`
	Meta string `json:"meta"`
}

// Binary is a struct that represents a binary file.
type Binary struct {
	Blob []byte `json:"blob"`
	Meta string `json:"meta"`
}

// AuthData is a struct that represents a hashed key, salt and tokens locally stored for user.
type AuthData struct {
	HashedKey    string `json:"hashed_key"`
	Salt         string `json:"salt"`
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
}
