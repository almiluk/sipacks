package v1

import "github.com/labstack/echo/v4"

// addPack godoc
// @Summary Add pack
// @Description Add new questions pack
// @Tags packs
// @Accept multipart/form-data
// @Produce json
// @Param pack formData file true "Pack data"
// @Success 200	{object} models.AddPackResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/packs [post]
func addPack(ctx echo.Context) error {
	return nil
}

// downloadPack godoc
// @Summary Download pack
// @Description Download questions pack
// @Tags packs
// @Produce octet-stream
// @Param id path int true "Pack ID"
// @Success 200
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /packs/{id} [get]
func downloadPack(ctx echo.Context) error {
	return nil
}

// listPacks godoc
// @Summary List packs
// @Description List packs with filters
// @Tags packs
// @Accept json
// @Produce json
// @Param filter body models.PackListRequest true "Filter"
// @Success 200 {object} models.PackListResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/packs [get]
func listPacks(ctx echo.Context) error {
	return nil
}
