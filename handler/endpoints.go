package handler

import (
	"fmt"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/google/uuid"
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
	resp.Id = uuid.New()
	return ctx.JSON(http.StatusOK, resp)
}
