package v1

import (
	"strings"

	"github.com/almiluk/sipacks/internal/controller/http/v1/models"
	"github.com/almiluk/sipacks/internal/entity"
	"github.com/labstack/echo/v4"
)

type SIPacksRouter struct {
	uc ISIPacksUC
}

func RegisterSIPacksRouter(handler *echo.Echo, uc ISIPacksUC) {
	router := SIPacksRouter{
		uc: uc,
	}

	handler.POST("/api/v1/packs", router.addPack)
	handler.GET("/api/v1/packs", router.listPacks)
	handler.GET("/packs/:guid", router.downloadPack)
}

// addPack godoc
// @Summary Add pack
// @Description Add new questions pack
// @Tags packs
// @Accept multipart/form-data
// @Produce json
// @Param pack formData file true "Pack data"
// @Success 200	{object} models.AddPackResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.PackResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/packs [post]
func (r SIPacksRouter) addPack(ctx echo.Context) error {
	packFile, err := ctx.FormFile("pack")
	if err != nil {
		return responseWithError(ctx, 400, "Request must contain pack file", err)
	}

	file, err := packFile.Open()
	if err != nil {
		return responseWithError(ctx, 500, "Can't open pack file", err)
	}
	defer file.Close()

	pack, err := r.uc.AddPack(ctx.Request().Context(), file, packFile.Size)
	if err != nil {
		return responseWithError(ctx, 500, "Can't add pack", err)
	}
	return ctx.JSON(200, pack)
}

// downloadPack godoc
// @Summary Download pack
// @Description Download questions pack
// @Tags packs
// @Produce octet-stream
// @Produce json
// @Param guid path string true "Pack guid" format(uuid)
// @Param filename query string false "Wanted package file name"
// @Success 200
// @Failure 404
// @Failure 500 {object} models.ErrorResponse
// @Router /packs/{guid} [get]
func (r SIPacksRouter) downloadPack(ctx echo.Context) error {
	guid := ctx.Param("guid")
	downloadFilename := ctx.QueryParam("filename")
	if downloadFilename == "" {
		downloadFilename = guid + ".siq"
	}

	storedFilename := r.uc.GetPackFilename(ctx.Request().Context(), guid)
	err := ctx.Attachment(storedFilename, downloadFilename)
	if err != nil {
		return err
	}

	err = r.uc.IncreaseDownloadsCounter(ctx.Request().Context(), guid)
	return err
}

// listPacks godoc
// @Summary List packs
// @Description List packs with filters
// @Tags packs
// @Accept json
// @Produce json
// @Param filter query models.PackListRequest true "Filter"
// @Success 200 {object} models.PackListResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/packs [get]
func (r SIPacksRouter) listPacks(ctx echo.Context) error {
	var filterBody models.PackListRequest
	if err := ctx.Bind(&filterBody); err != nil {
		return responseWithError(ctx, 400, "Incorrect filter", err)
	}

	filter := entity.PackFilter{
		Name:            filterBody.Name,
		Author:          filterBody.Author,
		MinCreationDate: filterBody.MinCreationDate,
		MaxCreationDate: filterBody.MaxCreationDate,
		SortBy:          filterBody.SortBy,
	}

	if filterBody.Tags != nil && *filterBody.Tags != "" {
		filter.Tags = strings.Split(*filterBody.Tags, ",")
	}

	packs, err := r.uc.GetPacks(ctx.Request().Context(), filter)
	if err != nil {
		return responseWithError(ctx, 500, "Cannot find packs", err)
	}

	packsResponse := make([]models.PackResponse, len(packs))
	for i := range packs {
		packsResponse[i] = models.PackResponse{
			GUID:         packs[i].GUID,
			Name:         packs[i].Name,
			Author:       packs[i].Author.Nickname,
			CreationDate: packs[i].CreationDate,
			FileSize:     packs[i].FileSize,
			DownloadsNum: packs[i].DownloadsNum,
			Tags:         make([]string, len(packs[i].Tags)),
		}
		for j := range packs[i].Tags {
			packsResponse[i].Tags[j] = packs[i].Tags[j].Name
		}
	}

	return ctx.JSON(200, models.PackListResponse{Packs: packsResponse, PacksNum: len(packs)})
}
