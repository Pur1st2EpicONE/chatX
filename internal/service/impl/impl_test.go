package impl

import (
	mockCache "chatX/internal/cache/mocks"
	"chatX/internal/config"
	"chatX/internal/errs"
	mockLogger "chatX/internal/logger/mocks"
	"chatX/internal/models"
	mockStorage "chatX/internal/repository/mocks"
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func newTestService(controller *gomock.Controller) (*Service, *mockLogger.MockLogger, *mockCache.MockCache, *mockStorage.MockStorage) {

	loggerMock := mockLogger.NewMockLogger(controller)
	cacheMock := mockCache.NewMockCache(controller)
	storageMock := mockStorage.NewMockStorage(controller)

	loggerMock.EXPECT().LogError(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().LogInfo(gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().LogWarn(gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().LogFatal(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	cfg := config.Service{
		MaxTitleLength:   200,
		MaxMessageLength: 1000,
		GetLimitDefault:  10,
		GetLimitMax:      100,
	}

	svc := NewService(loggerMock, cfg, cacheMock, storageMock)
	return svc, loggerMock, cacheMock, storageMock

}

func TestCreateChat_Success(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	svc, _, _, storageMock := newTestService(controller)

	chat := models.Chat{Title: "  test aboba  "}

	storageMock.EXPECT().CreateChat(gomock.Any(), gomock.AssignableToTypeOf(&models.Chat{})).Return(nil)

	res, err := svc.CreateChat(context.Background(), chat)
	assert.NoError(t, err)
	assert.Equal(t, "test aboba", res.Title)

}

func TestCreateChat_ValidateFail_TitleEmpty(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	svc, _, _, _ := newTestService(controller)

	chat := models.Chat{Title: "   "}

	_, err := svc.CreateChat(context.Background(), chat)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrTitleEmpty))

}

func TestCreateMessage_Success(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	svc, _, cacheMock, storageMock := newTestService(controller)

	msg := models.Message{ChatID: 1, Text: "  qwe  "}

	storageMock.EXPECT().CreateMessage(gomock.Any(), gomock.AssignableToTypeOf(&models.Message{})).Return(nil)
	cacheMock.EXPECT().Delete(msg.ChatID).Times(1)

	res, err := svc.CreateMessage(context.Background(), msg)
	assert.NoError(t, err)
	assert.Equal(t, "qwe", res.Text)
	assert.Equal(t, 1, res.ChatID)

}

func TestCreateMessage_ForeignKeyViolation_ReturnsChatNotFound(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	svc, _, cacheMock, storageMock := newTestService(controller)

	msg := models.Message{ChatID: 999, Text: "qweqweqwe"}

	pgErr := &pgconn.PgError{Code: pgerrcode.ForeignKeyViolation}

	storageMock.EXPECT().CreateMessage(gomock.Any(), gomock.AssignableToTypeOf(&models.Message{})).Return(pgErr)
	cacheMock.EXPECT().Delete(gomock.Any()).Times(0)

	_, err := svc.CreateMessage(context.Background(), msg)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrChatNotFound))

}

func TestDeleteChat_Success(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	svc, _, cacheMock, storageMock := newTestService(controller)

	chatID := 7

	storageMock.EXPECT().DeleteChat(gomock.Any(), chatID).Return(nil)
	cacheMock.EXPECT().Delete(chatID).Times(1)

	err := svc.DeleteChat(context.Background(), chatID)
	assert.NoError(t, err)

}

func TestDeleteChat_NotFound_ReturnsErrChatNotFound(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	svc, _, cacheMock, storageMock := newTestService(controller)

	chatID := 42

	storageMock.EXPECT().DeleteChat(gomock.Any(), chatID).Return(errs.ErrChatNotFound)
	cacheMock.EXPECT().Delete(gomock.Any()).Times(0)

	err := svc.DeleteChat(context.Background(), chatID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrChatNotFound))

}

func TestGetChat_CacheHit_NoStorageCall(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	svc, _, cacheMock, storageMock := newTestService(controller)

	chatID := 1
	chat := models.Chat{
		ID:    chatID,
		Title: "cached",
		Messages: []models.Message{
			{Text: "m1"},
			{Text: "m2"},
		},
	}

	cacheMock.EXPECT().Get(chatID).Return(chat, nil)
	storageMock.EXPECT().GetChat(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	res, err := svc.GetChat(context.Background(), chatID, "")
	assert.NoError(t, err)
	assert.Equal(t, "cached", res.Title)

}

func TestGetChat_CacheMiss_StorageSuccess_TrimsMessagesAndPutsToCache(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	loggerMock := mockLogger.NewMockLogger(controller)
	cacheMock := mockCache.NewMockCache(controller)
	storageMock := mockStorage.NewMockStorage(controller)
	loggerMock.EXPECT().LogError(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().LogInfo(gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().LogWarn(gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().LogFatal(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	cfg := config.Service{
		MaxTitleLength:   200,
		MaxMessageLength: 1000,
		GetLimitDefault:  10,
		GetLimitMax:      100,
	}

	svc := NewService(loggerMock, cfg, cacheMock, storageMock)

	chatID := 5
	chatFromDB := models.Chat{
		ID:    chatID,
		Title: "fromdb",
		Messages: []models.Message{
			{Text: "m1"},
			{Text: "m2"},
			{Text: "m3"},
		},
	}

	cacheMock.EXPECT().Get(chatID).Return(models.Chat{}, errors.New("cache miss"))
	storageMock.EXPECT().GetChat(gomock.Any(), chatID, cfg.GetLimitMax).Return(chatFromDB, nil)
	cacheMock.EXPECT().Put(chatID, chatFromDB).Times(1)

	res, err := svc.GetChat(context.Background(), chatID, "2")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res.Messages))
	assert.Equal(t, "m1", res.Messages[0].Text)
	assert.Equal(t, "m2", res.Messages[1].Text)

}

func TestGetChat_InvalidLimitString_ReturnsErrInvalidLimit(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	svc, _, _, _ := newTestService(controller)

	_, err := svc.GetChat(context.Background(), 1, "not int")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrInvalidLimit))

}

func TestValidateLimit_Boundaries(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	svc, _, _, _ := newTestService(controller)

	lim, err := svc.validateLimit("")
	assert.NoError(t, err)
	assert.Equal(t, svc.config.GetLimitDefault, lim)

	lim, err = svc.validateLimit("0")
	assert.NoError(t, err)
	assert.Equal(t, svc.config.GetLimitDefault, lim)

	_, err = svc.validateLimit("-1")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrLimitTooSmall))

	tooLarge := strconv.Itoa(svc.config.GetLimitMax + 1)
	_, err = svc.validateLimit(tooLarge)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrLimitTooLarge))

}

func TestCreateChat_StorageError_LogsAndReturnsError(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	loggerMock := mockLogger.NewMockLogger(controller)
	cacheMock := mockCache.NewMockCache(controller)
	storageMock := mockStorage.NewMockStorage(controller)

	loggerMock.EXPECT().LogError(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)
	loggerMock.EXPECT().LogInfo(gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().LogWarn(gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().LogFatal(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	cfg := config.Service{
		MaxTitleLength:   200,
		MaxMessageLength: 1000,
		GetLimitDefault:  10,
		GetLimitMax:      100,
	}
	svc := NewService(loggerMock, cfg, cacheMock, storageMock)

	chat := models.Chat{Title: "  asdadqwd  "}

	storageErr := errors.New("db down")
	storageMock.EXPECT().CreateChat(gomock.Any(), gomock.AssignableToTypeOf(&models.Chat{})).Return(storageErr)

	res, err := svc.CreateChat(context.Background(), chat)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db down")
	assert.Equal(t, models.Chat{}, res)

}

func TestCreateMessage_ValidationFails_ReturnsError(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	svc, _, cacheMock, storageMock := newTestService(controller)

	msg := models.Message{ChatID: 1, Text: "   "}

	cacheMock.EXPECT().Delete(gomock.Any()).Times(0)
	storageMock.EXPECT().CreateMessage(gomock.Any(), gomock.Any()).Times(0)

	res, err := svc.CreateMessage(context.Background(), msg)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrMessageEmpty))
	assert.Equal(t, models.Message{}, res)

}

func TestCreateMessage_StorageError_LogsAndReturnsError(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	loggerMock := mockLogger.NewMockLogger(controller)
	cacheMock := mockCache.NewMockCache(controller)
	storageMock := mockStorage.NewMockStorage(controller)

	loggerMock.EXPECT().LogError(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)
	loggerMock.EXPECT().LogInfo(gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().LogWarn(gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().LogFatal(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	cfg := config.Service{
		MaxTitleLength:   200,
		MaxMessageLength: 1000,
		GetLimitDefault:  10,
		GetLimitMax:      100,
	}
	svc := NewService(loggerMock, cfg, cacheMock, storageMock)

	message := models.Message{ChatID: 1, Text: "qwe"}

	storageErr := errors.New("db unavailable")
	storageMock.EXPECT().CreateMessage(gomock.Any(), gomock.AssignableToTypeOf(&models.Message{})).Return(storageErr)
	cacheMock.EXPECT().Delete(gomock.Any()).Times(0)

	res, err := svc.CreateMessage(context.Background(), message)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db unavailable")
	assert.Equal(t, models.Message{}, res)

}

func TestDeleteChat_StorageError_LogsAndReturnsError(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	loggerMock := mockLogger.NewMockLogger(controller)
	cacheMock := mockCache.NewMockCache(controller)
	storageMock := mockStorage.NewMockStorage(controller)

	storageErr := errors.New("db unavailable")
	loggerMock.EXPECT().LogError("service — failed to delete chat", storageErr, "chatID", 1, "layer", "service.impl").Times(1)
	loggerMock.EXPECT().LogInfo(gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().LogWarn(gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().LogFatal(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	svc := NewService(loggerMock, config.Service{}, cacheMock, storageMock)
	chatID := 1

	storageMock.EXPECT().DeleteChat(gomock.Any(), chatID).Return(storageErr)
	cacheMock.EXPECT().Delete(gomock.Any()).Times(0)

	err := svc.DeleteChat(context.Background(), chatID)

	assert.Error(t, err)
	assert.Equal(t, storageErr, err)

}

func TestGetChat_StorageError_LogsAndReturnsError(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	loggerMock := mockLogger.NewMockLogger(controller)
	cacheMock := mockCache.NewMockCache(controller)
	storageMock := mockStorage.NewMockStorage(controller)

	chatID := 1
	storageErr := errors.New("db unavailable")

	loggerMock.EXPECT().LogError("service — failed to get chat", storageErr, "chatID", chatID, "layer", "service.impl").Times(1)
	loggerMock.EXPECT().LogInfo(gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().LogWarn(gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().LogFatal(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	cfg := config.Service{GetLimitMax: 100}
	svc := NewService(loggerMock, cfg, cacheMock, storageMock)

	cacheMock.EXPECT().Get(chatID).Return(models.Chat{}, errors.New("cache miss"))
	storageMock.EXPECT().GetChat(gomock.Any(), chatID, cfg.GetLimitMax).Return(models.Chat{}, storageErr)
	cacheMock.EXPECT().Put(gomock.Any(), gomock.Any()).Times(0)

	res, err := svc.GetChat(context.Background(), chatID, "")

	assert.Error(t, err)
	assert.Equal(t, storageErr, err)
	assert.Equal(t, models.Chat{}, res)

}

func TestValidateChat_TitleTooLong_ReturnsError(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	svc, _, _, _ := newTestService(controller)

	longTitle := strings.Repeat("a", svc.config.MaxTitleLength+1)
	chat := &models.Chat{Title: longTitle}

	err := svc.validateChat(chat)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrTitleTooLong))

}

func TestValidateMessage_TextTooLong_ReturnsError(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	svc, _, _, _ := newTestService(controller)

	longText := strings.Repeat("x", svc.config.MaxMessageLength+1)
	err := svc.validateMessage(&models.Message{Text: longText})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrMessageTooLong))

}
