package user

// Service is the interface for the user service.
type Service interface {
	// CreateUser makes a request to the CreateUser RPC to create a new user and updates the config.
	CreateUser(username, password string) error

	// LogInUser attempts user authorization on the remote server and if it fails, attempts local authorization.
	LogInUser(username, password string) error

	// RefreshToken makes a request to the RefreshUserSession RPC to refresh the user's session and updates the config.
	RefreshToken() error
}
