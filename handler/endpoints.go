package handler

import (
	"fmt"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/models"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) GetHello(ctx echo.Context, params generated.GetHelloParams) error {
	var resp generated.HelloResponse
	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) PostEstate(ctx echo.Context) error {
	var resp generated.EstateResponse

	body := new(EstateRequest)
	if err := ctx.Bind(body); err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	if err := validator.New().Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	estate := models.NewEstate(body.Width, body.Length)

	// Persisting estate model to database through repository
	s.Repository.EstatePersist(ctx.Request().Context(), estate)

	resp.Id = estate.UUID
	return ctx.JSON(http.StatusOK, resp)
}
