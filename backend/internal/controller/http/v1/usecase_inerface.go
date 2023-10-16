package v1

import (
	"context"

	"github.com/almiluk/sipacks/internal/entity"
)

type ISIPacksUC interface {
	AddPack(ctx context.Context, fileReader entity.ReaderReadAt, fileSize int64) (entity.Pack, error)
	GetPacks(ctx context.Context, filter entity.PackFilter) ([]entity.Pack, error)
	GetPackFilename(ctx context.Context, guid string) string
	IncreaseDownloadsCounter(ctx context.Context, guid string) error
}
