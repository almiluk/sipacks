package models

import (
	"time"
)

type PackListRequest struct {
	Name            *string    `json:"name" example:"name" query:"name"`
	Author          *string    `json:"author" example:"author" query:"author"`
	Tags            *string    `json:"tags" example:"tag1,tag2" query:"tags"`
	MinCreationDate *time.Time `json:"min_creation_date" example:"01.01.1970" query:"min_creation_date"`
	MaxCreationDate *time.Time `json:"max_creation_date" example:"01.01.1970" query:"max_creation_date"`
	SortBy          *string    `json:"sort_by" example:"creation_date" enums:"creation_date,downloads_num" query:"sort_by"`
}
