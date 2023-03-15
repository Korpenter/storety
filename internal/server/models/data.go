package models

import "github.com/google/uuid"

// Data is the data model.
type Data struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	Name    string
	Type    string
	Content []byte
}

// DataInfo is the data info model.
type DataInfo struct {
	Name string
	Type string
}
