package handler

import (
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/models"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func (s *Server) PostEstate(ctx echo.Context) error {
	body := new(EstateRequest)
	if err := ctx.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := validator.New().Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	estate := models.NewEstate(body.Width, body.Length)

	estate, err := s.Repository.EstatePersist(ctx.Request().Context(), estate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error saving data")
	}

	return ctx.JSON(http.StatusOK, generated.EstateResponse{
		Id: estate.UUID,
	})
}
