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

// godoc PackFileXmlInfo
// PackFileXmlInfo is a struct for parsing the content.xml file describing the pack.
type PackFileXmlInfo struct {
	XMLName xml.Name `xml:"package"`
	Name    string   `xml:"name,attr"`
	Date    string   `xml:"date,attr"`
	Id      string   `xml:"id,attr"`
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

			packXmlInfo := PackFileXmlInfo{}
			if err := xml.Unmarshal(byteValue, &packXmlInfo); err != nil {
				return entity.Pack{}, err
			}
			res := entity.Pack{
				Name:     packXmlInfo.Name,
				FileSize: uint32(fileSize),
				GUID:     packXmlInfo.Id,
				Tags:     make([]entity.Tag, len(packXmlInfo.Tags.Tags)),
			}

			if len(packXmlInfo.Info.Authors.Authors) > 0 {
				res.Author = entity.Author{
					Nickname: packXmlInfo.Info.Authors.Authors[0],
				}
			}

			for i := range packXmlInfo.Tags.Tags {
				res.Tags[i] = entity.Tag{
					Name: packXmlInfo.Tags.Tags[i],
				}
			}

			res.CreationDate, err = time.Parse("02.01.2006", packXmlInfo.Date)

			return res, err
		}
	}

	return entity.Pack{}, errors.New("incorrect pack file structure: no content.xml file")
}
