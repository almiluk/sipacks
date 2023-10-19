package usecase

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"io"
	"time"

	"github.com/almiluk/sipacks/internal/entity"
)

const (
	packInfoFileName = "content.xml"
)

var (
	ErrNoContentXMLFile = errors.New("incorrect pack file structure: no content.xml file")
)

// godoc PackFileXMLInfo
// PackFileXMLInfo is a struct for parsing the content.xml file describing the pack.
type PackFileXMLInfo struct {
	XMLName xml.Name `xml:"package"`
	Name    string   `xml:"name,attr"`
	Date    string   `xml:"date,attr"`
	ID      string   `xml:"id,attr"`
	Info    struct {
		XMLName xml.Name `xml:"info"`
		Authors struct {
			XMLName xml.Name `xml:"authors"`
			Authors []string `xml:"author"`
		} `xml:"authors"`
	} `xml:"info"`
	Tags struct {
		XMLName xml.Name `xml:"tags"`
		Tags    []string `xml:"tag"`
	} `xml:"tags"`
}

// godoc GetPackFileInfo
// GetPackFileInfo parses the pack file (as a .zip archive)
// and returns the pack info extracted from the content.xml file inside.
func GetPackFileInfo(fileReader io.ReaderAt, fileSize int64) (entity.Pack, error) {
	zipReader, err := zip.NewReader(fileReader, fileSize)
	if err != nil {
		return entity.Pack{}, err
	}

	for _, file := range zipReader.File {
		if file.Name == packInfoFileName {
			f, err := file.Open()
			if err != nil {
				return entity.Pack{}, err
			}

			byteValue, err := io.ReadAll(f)
			if err != nil {
				return entity.Pack{}, err
			}

			packXMLInfo := PackFileXMLInfo{}
			if err = xml.Unmarshal(byteValue, &packXMLInfo); err != nil {
				return entity.Pack{}, err
			}

			res := entity.Pack{
				Name:     packXMLInfo.Name,
				FileSize: uint32(fileSize),
				GUID:     packXMLInfo.ID,
				Tags:     make([]entity.Tag, len(packXMLInfo.Tags.Tags)),
			}

			if len(packXMLInfo.Info.Authors.Authors) > 0 {
				res.Author = entity.Author{
					Nickname: packXMLInfo.Info.Authors.Authors[0],
				}
			}

			for i := range packXMLInfo.Tags.Tags {
				res.Tags[i] = entity.Tag{
					Name: packXMLInfo.Tags.Tags[i],
				}
			}

			res.CreationDate, err = time.Parse("02.01.2006", packXMLInfo.Date)

			return res, err
		}
	}

	return entity.Pack{}, ErrNoContentXMLFile
}
