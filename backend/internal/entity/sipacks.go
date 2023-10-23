package entity

import (
	"io"
	"time"
)

type Author struct {
	Id       uint32
	Nickname string
}

type Pack struct {
	Id           uint32    `json:"id"`
	GUID         string    `json:"guid"`
	Name         string    `json:"name"`
	Author       Author    `json:"author"`
	CreationDate time.Time `json:"creation_date"`
	FileSize     uint32    `json:"file_size"`
	DownloadsNum uint32    `json:"downloads_num"`
	Tags         []Tag     `json:"tags"`
}

type PackFilter struct {
	Name            *string    `json:"name"`
	Author          *string    `json:"author"`
	Tags            []string   `json:"tags"`
	MinCreationDate *time.Time `json:"min_creation_date"`
	MaxCreationDate *time.Time `json:"max_creation_date"`
	SortBy          *string    `json:"sort_by" enums:"creation_date,downloads_num"`
}

type Tag struct {
	Id   uint32
	Name string
}

type ReaderReadAt interface {
	io.Reader
	io.ReaderAt
}
