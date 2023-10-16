package entity

import (
	"errors"
)

var (
	ErrPackAlreadyExists = errors.New("pack with the same GUID already exists")
)
