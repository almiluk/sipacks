package usecase

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/almiluk/sipacks/internal/entity"
)

type SIPacksUC struct {
	repo               IRepo
	absFileStoragePath string
}

func New(repo IRepo, fileStoragePath string) (*SIPacksUC, error) {
	absFileStoragePath, err := filepath.Abs(fileStoragePath)
	if err != nil {
		return nil, err
	}

	if err = os.Mkdir(absFileStoragePath, 0777); err != nil && !errors.Is(err, os.ErrExist) {
		return nil, err
	}

	return &SIPacksUC{
		repo:               repo,
		absFileStoragePath: absFileStoragePath,
	}, nil
}
func (uc *SIPacksUC) AddPack(ctx context.Context, fileReader entity.ReaderReadAt, fileSize int64) (entity.Pack, error) {
	// Extract pack info from the file
	pack, err := GetPackFileInfo(fileReader, fileSize)
	if err != nil {
		return entity.Pack{}, err
	}

	// Add pack metadata to the database
	err = uc.repo.AddPack(ctx, &pack)
	if errors.Is(err, entity.ErrPackAlreadyExists) {
		return pack, err
	} else if err != nil {
		return entity.Pack{}, err
	}

	// TODO: revert changes if file saving goes wrong

	// Add pack file to the storage
	filepath := uc.GetPackFilename(ctx, pack.GUID)
	file, err := os.Create(filepath)
	if err != nil {
		return entity.Pack{}, err
	}
	defer file.Close()

	_, err = io.Copy(file, fileReader)
	if err != nil {
		return entity.Pack{}, err
	}

	return pack, nil
}

func (uc *SIPacksUC) GetPacks(ctx context.Context, filter entity.PackFilter) ([]entity.Pack, error) {
	packs, err := uc.repo.GetPacks(ctx, filter)
	if err != nil {
		return nil, err
	}

	return packs, nil
}

func (us *SIPacksUC) GetPackFilename(ctx context.Context, guid string) string {
	filename := filepath.Join(us.absFileStoragePath, guid+".siq")
	return filename
}

func (us *SIPacksUC) IncreaseDownloadsCounter(ctx context.Context, guid string) error {
	return us.repo.IncreaseDownloadsCounter(ctx, guid)
}
