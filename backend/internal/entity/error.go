package entity

import (
	"errors"
)

var (
	ErrPackAlreadyExists = errors.New("pack with the same GUID already exists")
	ErrNoContentXMLFile  = errors.New("incorrect pack file structure: no content.xml file")
)
