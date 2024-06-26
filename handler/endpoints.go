package handler

import (
	"net/http"
	"sort"

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

	return ctx.JSON(http.StatusCreated, generated.EstateResponse{
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
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if oldTree != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "tree already exist in that coordinate")
	}
	// Done Check if the tree with the same coordinate already exists

	// Get Existing trees to calculate median
	trees, err := s.Repository.GetTreesByEstate(context, estate.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Create New Tree Entity
	newTree, err := models.NewTree(estate, uint16(body.X), uint16(body.Y), uint8(body.Height))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	treeValues := *trees
	treeValues = append(treeValues, *newTree)

	sort.Slice(treeValues, func(i, j int) bool {
		return treeValues[i].Height < treeValues[j].Height
	})

	estate.CalculateEstateTreeMedian(&treeValues)
	newTree.Estate = estate

	// Save New Tree Entity
	err = s.Repository.SaveTree(context, newTree)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, generated.EstateResponse{
		Id: newTree.UUID,
	})
}

func (s *Server) GetEstateIdStats(ctx echo.Context, id generated.EstateIDPathParam) error {
	context := ctx.Request().Context()

	// Start Check if the estate exist
	estate, err := s.Repository.GetEstate(context, id.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if estate == nil {
		return echo.NewHTTPError(http.StatusNotFound, "estate not found")
	}
	// Done Check if the estate exist

	return ctx.JSON(http.StatusOK, generated.EstateStatsResponse{
		Count:  int(estate.TreeCount),
		Max:    int(estate.MaxTreeHeight),
		Min:    int(estate.MinTreeHeight),
		Median: int(estate.MedianTreeHeight),
	})
}

func (s *Server) GetEstateIdDronePlan(ctx echo.Context, id generated.EstateIDPathParam, params generated.GetEstateIdDronePlanParams) error {
	context := ctx.Request().Context()
	var maxDistance *uint16
	if params.MaxDistance != nil {
		maxDistanceRaw := uint16(*params.MaxDistance)
		maxDistance = &maxDistanceRaw
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

	trees, err := s.Repository.GetTreesByEstate(context, estate.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	drone := models.NewDrone(estate, trees, maxDistance)
	drone.StartFlight()

	if maxDistance == nil {
		return ctx.JSON(http.StatusOK, generated.DronePlanResponse{
			Distance: int(drone.Travelled),
		})
	} else {
		lastCoordinateX := int(drone.LastCoordinateX)
		lastCoordinateY := int(drone.LastCoordinateY)

		return ctx.JSON(http.StatusOK, generated.DronePlanResponse{
			Distance: int(drone.Travelled),
			Rest: &generated.DroneRestResponse{
				X: &lastCoordinateX,
				Y: &lastCoordinateY,
			},
		})
	}
}
