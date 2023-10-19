package usecase_test

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/almiluk/sipacks/internal/entity"
	"github.com/almiluk/sipacks/internal/usecase"
)

func TestGetPackFileInfo(t *testing.T) {
	type Usecase struct {
		name     string
		filename string
		expected entity.Pack
		err      error
	}

	testCases := []Usecase{
		{
			name:     "test_pack.siq",
			filename: "pack_examples/test_pack.siq",
			expected: entity.Pack{
				Name: "Сомнительная смесь #1",
				Author: entity.Author{
					Nickname: "brave_new_rave",
				},
				GUID:         "828343ff-82ba-4ef2-ad5b-f3d1e0a3c2b3",
				CreationDate: time.Date(2023, 10, 12, 0, 0, 0, 0, time.UTC),
				FileSize:     46199834,
				Tags: []entity.Tag{
					{
						Name: "Фильмы",
					},
					{
						Name: "Музыка",
					},
					{
						Name: "Нейронки",
					},
					{
						Name: "Актёры",
					},
					{
						Name: "Еда",
					},
					{
						Name: "Сериалы",
					},
					{
						Name: "Игры",
					},
				},
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		file, err := os.Open(tc.filename)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		fileStat, err := file.Stat()
		if err != nil {
			t.Fatal(err)
		}

		fileInfo, err := usecase.GetPackFileInfo(file, fileStat.Size())
		if errors.Is(err, tc.err) {
			t.Fatalf("expected not nil error: %v, got: %v", tc.err, err)
		}

		if !ComparePacks(fileInfo, tc.expected) {
			t.Fatalf("expected: %+v\ngot: %+v", tc.expected, fileInfo)
		}

		fmt.Printf("%+v\n", fileInfo)
	}

}

func ComparePacks(p1, p2 entity.Pack) bool {
	same := p1.ID == p2.ID &&
		p1.Name == p2.Name &&
		p1.Author.ID == p2.Author.ID &&
		p1.Author.Nickname == p2.Author.Nickname &&
		p1.CreationDate == p2.CreationDate &&
		p1.FileSize == p2.FileSize &&
		p1.DownloadsNum == p2.DownloadsNum &&
		len(p1.Tags) == len(p2.Tags)

	if same {
		for i, tag := range p1.Tags {
			if tag != p2.Tags[i] {
				return false
			}
		}
		return true
	}

	return same
}
