package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
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

		assert.Equal(t, http.StatusOK, rec.Code)
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
