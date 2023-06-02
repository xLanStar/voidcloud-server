package util

import (
	"nullprogram.com/x/uuid"
)

var uuidGenerator = uuid.NewGen()

func NewUUID() string {
	return uuidGenerator.NewV4().String()
}
