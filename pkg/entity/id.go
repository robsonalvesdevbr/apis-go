package entity

import "github.com/google/uuid"

type ID = uuid.UUID

func NewID() (ID, error) {
	return uuid.NewV7()
}
