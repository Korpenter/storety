package helpers

import "log"

// LogError logs an error if it is not nil.
func LogError(err error) error {
	if err != nil {
		log.Println("Error running command: ", err)
	}
	return err
}
