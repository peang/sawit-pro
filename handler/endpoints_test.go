package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/models"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPostEstate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	s := &Server{
		Repository: mockRepo,
	}

	requestBody := `{"width": 10, "length": 20}`

	req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewBufferString(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	mockRepo.EXPECT().SaveEstate(c.Request().Context(), gomock.Any()).Return(nil)

	if assert.NoError(t, s.PostEstate(c)) {
		var responseBody generated.EstateResponse
		err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
		if err != nil {
			t.Fatalf("failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.NotEmpty(t, responseBody.Id, "id should not be empty")
	}
}

func TestPostEstate_BadPayload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	s := &Server{
		Repository: mockRepo,
	}

	requestBody := `{"invalid_json": "missing required fields"}`
	req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewBufferString(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := s.PostEstate(c)

	if httpErr, ok := err.(*echo.HTTPError); ok {
		statusCode := httpErr.Code

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	}
}

func TestPostEstate_ErrorPersisting(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	s := &Server{
		Repository: mockRepo,
	}

	requestBody := `{"width": 10, "length": 20}`

	req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewBufferString(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	mockRepo.EXPECT().SaveEstate(c.Request().Context(), gomock.Any()).Return(errors.New("error"))

	err := s.PostEstate(c)
	if httpErr, ok := err.(*echo.HTTPError); ok {
		statusCode := httpErr.Code

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, statusCode, "error saving data")
	}
}

func TestPostTree(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	estateId := uint64(1)
	estateUuid := uuid.New()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	mockEstate := models.Estate{
		ID:     estateId,
		UUID:   estateUuid.String(),
		Width:  10,
		Length: 10,
	}

	s := &Server{
		Repository: mockRepo,
	}

	requestBody := `{"x": 1, "y": 1, "height": 10}`

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/estate/%s/tree", estateUuid), bytes.NewBufferString(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	mockRepo.EXPECT().GetEstate(c.Request().Context(), estateUuid.String()).Return(&mockEstate, nil)
	mockRepo.EXPECT().GetTreeByCoordinate(c.Request().Context(), estateId, uint16(1), uint16(1)).Return(nil, nil)
	mockRepo.EXPECT().SaveTree(c.Request().Context(), gomock.Any())

	if assert.NoError(t, s.PostEstateIdTree(c, estateUuid)) {
		var responseBody generated.TreeResponse
		err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
		if err != nil {
			t.Fatalf("failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.NotEmpty(t, responseBody.Id, "id should not be empty")
	}
}

func TestPostTree_ErrorGetEstate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	estateUuid := uuid.New()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	s := &Server{
		Repository: mockRepo,
	}

	requestBody := `{"x": 1, "y": 1, "height": 10}`

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/estate/%s/tree", estateUuid), bytes.NewBufferString(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	mockRepo.EXPECT().GetEstate(c.Request().Context(), estateUuid.String()).Return(nil, errors.New("error"))

	err := s.PostEstateIdTree(c, estateUuid)
	if httpErr, ok := err.(*echo.HTTPError); ok {
		statusCode := httpErr.Code

		assert.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	}
}

func TestPostTree_ErrorEstateNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	estateUuid := uuid.New()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	s := &Server{
		Repository: mockRepo,
	}

	requestBody := `{"x": 1, "y": 1, "height": 10}`

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/estate/%s/tree", estateUuid), bytes.NewBufferString(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	mockRepo.EXPECT().GetEstate(c.Request().Context(), estateUuid.String()).Return(nil, nil)

	err := s.PostEstateIdTree(c, estateUuid)
	if httpErr, ok := err.(*echo.HTTPError); ok {
		statusCode := httpErr.Code
		statusMessage := httpErr.Message

		assert.Error(t, err)
		assert.Equal(t, http.StatusNotFound, statusCode)
		assert.Equal(t, "estate not found", statusMessage)
	}
}

func TestPostTree_TreeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	estateId := uint64(1)
	estateUuid := uuid.New()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	mockEstate := models.Estate{
		ID:     estateId,
		UUID:   estateUuid.String(),
		Width:  10,
		Length: 10,
	}
	s := &Server{
		Repository: mockRepo,
	}

	requestBody := `{"x": 1, "y": 1, "height": 10}`

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/estate/%s/tree", estateUuid), bytes.NewBufferString(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	mockRepo.EXPECT().GetEstate(c.Request().Context(), estateUuid.String()).Return(&mockEstate, nil)
	mockRepo.EXPECT().GetTreeByCoordinate(c.Request().Context(), estateId, uint16(1), uint16(1)).Return(nil, errors.New("error"))

	err := s.PostEstateIdTree(c, estateUuid)
	if httpErr, ok := err.(*echo.HTTPError); ok {
		statusCode := httpErr.Code
		statusMessage := httpErr.Message

		assert.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
		assert.Equal(t, "error", statusMessage)
	}
}

func TestPostTree_TreeAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	estateId := uint64(1)
	treeId := uint64(1)
	estateUuid := uuid.New()
	treeUuid := uuid.New()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	mockEstate := models.Estate{
		ID:     estateId,
		UUID:   estateUuid.String(),
		Width:  10,
		Length: 10,
	}
	mockTree := models.Tree{
		ID:     treeId,
		UUID:   treeUuid.String(),
		X:      1,
		Y:      1,
		Height: 10,
	}

	s := &Server{
		Repository: mockRepo,
	}

	requestBody := `{"x": 1, "y": 1, "height": 10}`

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/estate/%s/tree", estateUuid), bytes.NewBufferString(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	mockRepo.EXPECT().GetEstate(c.Request().Context(), estateUuid.String()).Return(&mockEstate, nil)
	mockRepo.EXPECT().GetTreeByCoordinate(c.Request().Context(), estateId, uint16(1), uint16(1)).Return(&mockTree, nil)

	err := s.PostEstateIdTree(c, estateUuid)
	if httpErr, ok := err.(*echo.HTTPError); ok {
		statusCode := httpErr.Code
		statusMessage := httpErr.Message

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Equal(t, "tree already exist in that coordinate", statusMessage)
	}
}

func TestPostTree_TreeErrorWhenCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	estateId := uint64(1)
	estateUuid := uuid.New()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	mockEstate := models.Estate{
		ID:     estateId,
		UUID:   estateUuid.String(),
		Width:  10,
		Length: 10,
	}
	s := &Server{
		Repository: mockRepo,
	}

	requestBody := `{"x": 1, "y": 1, "height": 10}`

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/estate/%s/tree", estateUuid), bytes.NewBufferString(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	mockRepo.EXPECT().GetEstate(c.Request().Context(), estateUuid.String()).Return(&mockEstate, nil)
	mockRepo.EXPECT().GetTreeByCoordinate(c.Request().Context(), estateId, uint16(1), uint16(1)).Return(nil, nil)
	mockRepo.EXPECT().SaveTree(c.Request().Context(), gomock.Any()).Return(errors.New("error"))

	err := s.PostEstateIdTree(c, estateUuid)
	if httpErr, ok := err.(*echo.HTTPError); ok {
		statusCode := httpErr.Code
		statusMessage := httpErr.Message

		assert.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
		assert.Equal(t, "error", statusMessage)
	}
}

func TestGetEstateStats_ErrorWhenGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	s := &Server{
		Repository: mockRepo,
	}

	estateUuid := uuid.New()

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/estate/%s/stats", estateUuid.String()), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	mockRepo.EXPECT().GetEstate(c.Request().Context(), gomock.Any()).Return(nil, errors.New("error"))

	err := s.GetEstateIdStats(c, estateUuid)
	if httpErr, ok := err.(*echo.HTTPError); ok {
		statusCode := httpErr.Code
		statusMessage := httpErr.Message

		assert.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
		assert.Equal(t, "error", statusMessage)
	}
}

func TestGetEstateStats_EstateNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	s := &Server{
		Repository: mockRepo,
	}

	estateUuid := uuid.New()

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/estate/%s/stats", estateUuid.String()), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	mockRepo.EXPECT().GetEstate(c.Request().Context(), gomock.Any()).Return(nil, nil)

	err := s.GetEstateIdStats(c, estateUuid)
	if httpErr, ok := err.(*echo.HTTPError); ok {
		statusCode := httpErr.Code
		statusMessage := httpErr.Message

		assert.Error(t, err)
		assert.Equal(t, http.StatusNotFound, statusCode)
		assert.Equal(t, "estate not found", statusMessage)
	}
}

func TestGetEstateStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	estateId := uint64(1)
	estateUuid := uuid.New()
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	mockEstate := models.Estate{
		ID:               estateId,
		UUID:             estateUuid.String(),
		Width:            10,
		Length:           10,
		TreeCount:        2,
		MaxTreeHeight:    10,
		MinTreeHeight:    1,
		MedianTreeHeight: 5,
	}

	mockResponses := generated.EstateStatsResponse{
		Count:  2,
		Max:    10,
		Min:    1,
		Median: 5,
	}

	s := &Server{
		Repository: mockRepo,
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/estate/%s/stats", estateUuid.String()), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	mockRepo.EXPECT().GetEstate(c.Request().Context(), estateUuid.String()).Return(&mockEstate, nil)

	if assert.NoError(t, s.GetEstateIdStats(c, estateUuid)) {
		var responseBody generated.EstateStatsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
		if err != nil {
			t.Fatalf("failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, mockResponses, responseBody)
	}
}
