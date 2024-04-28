package handler

import (
	"fmt"
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

	err := s.Repository.SaveEstate(ctx.Request().Context(), estate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error saving data")
	}

	return ctx.JSON(http.StatusOK, generated.EstateResponse{
		Id: estate.UUID,
	})
}

func (s *Server) PostEstateIdTree(ctx echo.Context, id generated.EstateIDPathParam) error {
	context := ctx.Request().Context()
	body := new(TreeRequest)
	if err := ctx.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := validator.New().Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Start Check if the estate exist
	estate, err := s.Repository.GetEstate(context, id.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if estate == nil {
		return echo.NewHTTPError(http.StatusNotFound, "estate not found")
	}
	// Done Check if the estate exist

	// Start Check if the tree with the same coordinate already exists
	oldTree, err := s.Repository.GetTreeByCoordinate(context, estate.ID, uint16(body.X), uint16(body.Y))
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if oldTree != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "tree already exist in that coordinate")
	}
	// Done Check if the tree with the same coordinate already exists

	// Create New Tree Entity
	newTree, err := models.NewTree(estate, uint16(body.X), uint16(body.Y), uint8(body.Height))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Save New Tree Entity
	err = s.Repository.SaveTree(context, newTree)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, generated.EstateResponse{
		Id: newTree.UUID,
	})
}
