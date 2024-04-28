package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SawitProRecruitment/UserService/models"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPostEstate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	mockEstate := models.NewEstate(10, 10)

	s := &Server{
		Repository: mockRepo,
	}

	requestBody := `{"width": 10, "length": 20}`

	req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewBufferString(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	mockRepo.EXPECT().EstatePersist(c.Request().Context(), gomock.Any()).Return(mockEstate, nil)

	if assert.NoError(t, s.PostEstate(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, fmt.Sprintf(`{"id":"%s"}`, mockEstate.UUID), strings.TrimSpace(rec.Body.String()))
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

	mockRepo.EXPECT().EstatePersist(c.Request().Context(), gomock.Any()).Return(nil, errors.New("error"))

	err := s.PostEstate(c)
	if httpErr, ok := err.(*echo.HTTPError); ok {
		statusCode := httpErr.Code

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, statusCode, "error saving data")
	}
}
