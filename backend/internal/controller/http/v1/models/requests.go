package models

import (
	"os"
	"time"
)

type AddPackRequest struct {
	File os.File
}

type PackListRequest struct {
	Name            string    `json:"name" example:"name"`
	Author          string    `json:"author" example:"author"`
	Tags            []string  `json:"tags" example:"tags"`
	MinCreationDate time.Time `json:"min_creation_date" example:"01.01.1970"`
	MaxCreationDate time.Time `json:"max_creation_date" example:"01.01.1970"`
	SortBy          string    `json:"sort_by" example:"creation_date" enums:"creation_date,downloads_num"`
}
