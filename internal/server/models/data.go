package models

import "github.com/google/uuid"

type Data struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	Name    string
	Type    string
	Content []byte
}

type DataInfo struct {
	Name string
	Type string
}
