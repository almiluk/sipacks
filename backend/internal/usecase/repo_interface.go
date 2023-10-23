package usecase

import (
	"context"

	"github.com/almiluk/sipacks/internal/entity"
)

type IRepo interface {
	AddPack(ctx context.Context, pack *entity.Pack) error
	GetPacks(ctx context.Context, filter entity.PackFilter) ([]entity.Pack, error)
	IncreaseDownloadsCounter(ctx context.Context, guid string) error
}
