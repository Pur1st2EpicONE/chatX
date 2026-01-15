package v1

import (
	"chatX/internal/errs"
	"chatX/internal/models"
	"chatX/internal/service/mocks"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func setupRouter(h *Handler) *gin.Engine {

	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/chats", h.CreateChat)
	router.POST("/chats/:id/messages", h.CreateMessage)
	router.DELETE("/chats/:id", h.DeleteChat)
	router.GET("/chats/:id", h.GetChat)

	return router

}

func TestHandler_CreateChat_OK(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	service := mocks.NewMockService(controller)
	handler := NewHandler(service)
	router := setupRouter(handler)

	service.EXPECT().CreateChat(gomock.Any(), models.Chat{Title: "qweqwe"}).Return(models.Chat{ID: 1, Title: "qweqwe", CreatedAt: time.Now()}, nil)

	req := httptest.NewRequest(http.MethodPost, "/chats", strings.NewReader(`{"title":"qweqwe"}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"id":1`)
	assert.Contains(t, w.Body.String(), `"title":"qweqwe"`)

}

func TestHandler_CreateChat_InvalidJSON(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	service := mocks.NewMockService(controller)
	handler := NewHandler(service)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/chats", strings.NewReader("{invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), errs.ErrInvalidJSON.Error())

}

func TestHandler_CreateMessage_OK(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	service := mocks.NewMockService(controller)
	handler := NewHandler(service)
	router := setupRouter(handler)

	service.EXPECT().CreateMessage(gomock.Any(), models.Message{ChatID: 1, Text: "aboba"}).
		Return(models.Message{ID: 10, ChatID: 1, Text: "aboba", CreatedAt: time.Now()}, nil)

	req := httptest.NewRequest(http.MethodPost, "/chats/1/messages", strings.NewReader(`{"text":"aboba"}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"chat_id":1`)

}

func TestHandler_CreateMessage_InvalidChatID(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	service := mocks.NewMockService(controller)
	handler := NewHandler(service)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/chats/abc/messages", strings.NewReader(`{"text":"aboba"}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), errs.ErrInvalidChatID.Error())

}

func TestHandler_DeleteChat_OK(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	service := mocks.NewMockService(controller)
	handler := NewHandler(service)
	router := setupRouter(handler)

	service.EXPECT().DeleteChat(gomock.Any(), 1).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/chats/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

}

func TestHandler_DeleteChat_NotFound(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	service := mocks.NewMockService(controller)
	handler := NewHandler(service)
	router := setupRouter(handler)

	service.EXPECT().DeleteChat(gomock.Any(), 1).Return(errs.ErrChatNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/chats/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

}

func TestHandler_GetChat_OK(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	service := mocks.NewMockService(controller)
	handler := NewHandler(service)
	router := setupRouter(handler)

	service.EXPECT().GetChat(gomock.Any(), 1, "").
		Return(models.Chat{ID: 1, Title: "chat", Messages: []models.Message{{ID: 1, ChatID: 1, Text: "aboba"}}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/chats/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"messages"`)

}

func TestCreateChat_ServiceError(t *testing.T) {

	gin.SetMode(gin.TestMode)

	controller := gomock.NewController(t)
	defer controller.Finish()

	svc := mocks.NewMockService(controller)
	h := NewHandler(svc)

	svc.EXPECT().CreateChat(gomock.Any(), gomock.Any()).Return(models.Chat{}, errs.ErrTitleEmpty)

	router := gin.New()
	router.POST("/chats", h.CreateChat)

	body := `{"title":"   "}`
	req := httptest.NewRequest(http.MethodPost, "/chats", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

}

func TestHandler_CreateMessage_InvalidJSON(t *testing.T) {

	gin.SetMode(gin.TestMode)

	controller := gomock.NewController(t)
	defer controller.Finish()

	svc := mocks.NewMockService(controller)
	h := NewHandler(svc)

	router := gin.New()
	router.POST("/chats/:id/messages", h.CreateMessage)

	req := httptest.NewRequest(http.MethodPost, "/chats/1/messages", strings.NewReader(`{invalid json`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), errs.ErrInvalidJSON.Error())

}

func TestHandler_CreateMessage_ServiceError(t *testing.T) {

	gin.SetMode(gin.TestMode)

	controller := gomock.NewController(t)
	defer controller.Finish()

	svc := mocks.NewMockService(controller)
	h := NewHandler(svc)

	expectedErr := errs.ErrMessageEmpty

	svc.EXPECT().CreateMessage(gomock.Any(), models.Message{ChatID: 1, Text: ""}).Return(models.Message{}, expectedErr)

	router := gin.New()
	router.POST("/chats/:id/messages", h.CreateMessage)

	req := httptest.NewRequest(http.MethodPost, "/chats/1/messages", strings.NewReader(`{"text":""}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), expectedErr.Error())

}

func TestHandler_DeleteChat_InvalidChatID(t *testing.T) {

	gin.SetMode(gin.TestMode)

	controller := gomock.NewController(t)
	defer controller.Finish()

	svc := mocks.NewMockService(controller)
	h := NewHandler(svc)

	router := gin.New()
	router.DELETE("/chats/:id", h.DeleteChat)

	req := httptest.NewRequest(http.MethodDelete, "/chats/abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), errs.ErrInvalidChatID.Error())

}

func TestHandler_GetChat_InvalidChatID(t *testing.T) {

	gin.SetMode(gin.TestMode)

	controller := gomock.NewController(t)
	defer controller.Finish()

	svc := mocks.NewMockService(controller)
	h := NewHandler(svc)

	router := gin.New()
	router.GET("/chats/:id", h.GetChat)

	req := httptest.NewRequest(http.MethodGet, "/chats/xyz", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), errs.ErrInvalidChatID.Error())

}

func TestHandler_GetChat_NotFound(t *testing.T) {

	gin.SetMode(gin.TestMode)

	controller := gomock.NewController(t)
	defer controller.Finish()

	svc := mocks.NewMockService(controller)
	h := NewHandler(svc)

	svc.EXPECT().GetChat(gomock.Any(), 1, "").Return(models.Chat{}, errs.ErrChatNotFound)

	router := gin.New()
	router.GET("/chats/:id", h.GetChat)

	req := httptest.NewRequest(http.MethodGet, "/chats/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
	require.Contains(t, w.Body.String(), errs.ErrChatNotFound.Error())

}

func TestHandler_GetChat_InternalError(t *testing.T) {

	gin.SetMode(gin.TestMode)

	controller := gomock.NewController(t)
	defer controller.Finish()

	svc := mocks.NewMockService(controller)
	h := NewHandler(svc)

	svc.EXPECT().GetChat(gomock.Any(), 1, "").Return(models.Chat{}, errors.New("db is down"))

	router := gin.New()
	router.GET("/chats/:id", h.GetChat)

	req := httptest.NewRequest(http.MethodGet, "/chats/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.Contains(t, w.Body.String(), errs.ErrInternal.Error())

}
