package models

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Meta     string `json:"meta"`
}

type Card struct {
	Number  string `json:"number"`
	Expires string `json:"expires"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	CVV     string `json:"cvv"`
	Meta    string `json:"meta"`
}

type Text struct {
	Text string `json:"text"`
	Meta string `json:"meta"`
}

type Binary struct {
	Blob []byte `json:"blob"`
	Meta string `json:"meta"`
}
