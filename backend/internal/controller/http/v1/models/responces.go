package models

import "time"

type ErrorResponse struct {
	Message string `json:"message" example:"message"`
	Error   string `json:"error" example:"error"`
}

type AddPackResponse struct {
	Name         string    `json:"name" example:"name"`
	Author       string    `json:"author" example:"author"`
	CreationDate time.Time `json:"creation_date" example:"creation_date"`
}

type PackListResponse struct {
	Packs    []PackResponse `json:"packs"`
	PacksNum uint32         `json:"packs_num" example:"0"`
}

type PackResponse struct {
	GUID         string    `json:"guid" example:"00000000-0000-0000-0000-000000000000"`
	Name         string    `json:"name" example:"name"`
	Author       string    `json:"author" example:"author"`
	CreationDate time.Time `json:"creation_date" example:"creation_date"`
	FileSize     uint32    `json:"file_size" example:"0"`
	DownloadsNum uint32    `json:"downloads_num" example:"0"`
}
