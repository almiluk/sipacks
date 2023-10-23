package usecase

type PacksUC struct {
	repo IRepo
}

func New(repo IRepo) *PacksUC {
	return &PacksUC{
		repo: repo,
	}
}
