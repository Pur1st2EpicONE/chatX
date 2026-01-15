package postgres_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"chatX/internal/config"
	"chatX/internal/errs"
	"chatX/internal/logger"
	"chatX/internal/models"
	"chatX/internal/repository/postgres"

	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testStorage *postgres.Storage

func TestMain(m *testing.M) {

	logger, _ := logger.NewLogger(config.Logger{Debug: true})

	err := godotenv.Load("../../../.env")
	if err != nil {
		logger.LogFatal("Error loading .env file: %v", err)
	}

	cfg := config.Storage{
		Host:     "localhost",
		Port:     "5434",
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   "chatX_test",
		SSLMode:  "disable",
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.LogFatal("failed to connect to test DB: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.LogFatal("failed to get sql.DB from gorm: %v", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		logger.LogFatal("failed to set goose dialect: %v", err)
	}

	if err := goose.Up(sqlDB, "../../../migrations"); err != nil {
		logger.LogFatal("failed to apply migrations: %v", err)
	}

	testStorage = postgres.NewStorage(logger, cfg, db)

	exitCode := m.Run()

	testStorage.Close()
	os.Exit(exitCode)

}

func TestChatLifecycle(t *testing.T) {

	ctx := context.Background()

	chat := &models.Chat{Title: "Integration Chat", CreatedAt: time.Now().UTC()}

	err := testStorage.CreateChat(ctx, chat)
	if err != nil {
		t.Fatalf("CreateChat failed: %v", err)
	}

	if chat.ID == 0 {
		t.Fatal("chat ID not set after creation")
	}

	msg := &models.Message{
		ChatID:    chat.ID,
		Text:      "Hello world",
		CreatedAt: time.Now().UTC(),
	}

	err = testStorage.CreateMessage(ctx, msg)
	if err != nil {
		t.Fatalf("CreateMessage failed: %v", err)
	}

	if msg.ID == 0 {
		t.Fatal("message ID not set after creation")
	}

	gotChat, err := testStorage.GetChat(ctx, chat.ID, 10)
	if err != nil {
		t.Fatalf("GetChat failed: %v", err)
	}

	if gotChat.ID != chat.ID || gotChat.Title != chat.Title {
		t.Fatal("GetChat returned wrong chat data")
	}

	if len(gotChat.Messages) != 1 || gotChat.Messages[0].Text != msg.Text {
		t.Fatal("GetChat returned wrong messages")
	}

	err = testStorage.DeleteChat(ctx, chat.ID)
	if err != nil {
		t.Fatalf("DeleteChat failed: %v", err)
	}

	_, err = testStorage.GetChat(ctx, chat.ID, 10)
	if err == nil || !errors.Is(err, errs.ErrChatNotFound) {
		t.Fatal("expected ErrChatNotFound after delete")
	}

}

func TestDeleteNonExistingChat(t *testing.T) {
	ctx := context.Background()
	err := testStorage.DeleteChat(ctx, 9999)
	if err == nil || !errors.Is(err, errs.ErrChatNotFound) {
		t.Fatalf("expected ErrChatNotFound, got %v", err)
	}
}

func TestGetChatWithLimit(t *testing.T) {

	ctx := context.Background()

	chat := &models.Chat{Title: "Limit Chat", CreatedAt: time.Now().UTC()}

	err := testStorage.CreateChat(ctx, chat)
	if err != nil {
		t.Fatalf("CreateChat failed: %v", err)
	}

	for i := 1; i <= 5; i++ {

		msg := &models.Message{
			ChatID:    chat.ID,
			Text:      fmt.Sprintf("Message %d", i),
			CreatedAt: time.Now().Add(time.Duration(i) * time.Second).UTC(),
		}

		err := testStorage.CreateMessage(ctx, msg)
		if err != nil {
			t.Fatalf("CreateMessage failed: %v", err)
		}

	}

	gotChat, err := testStorage.GetChat(ctx, chat.ID, 3)
	if err != nil {
		t.Fatalf("GetChat failed: %v", err)
	}

	if len(gotChat.Messages) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(gotChat.Messages))
	}

	prev := gotChat.Messages[0].CreatedAt
	for _, message := range gotChat.Messages[1:] {
		if !prev.After(message.CreatedAt) && !prev.Equal(message.CreatedAt) {
			t.Fatalf("messages not in descending order")
		}
		prev = message.CreatedAt
	}

}

func TestStorageClose(t *testing.T) {
	if testStorage == nil {
		t.Fatal("testStorage is nil")
	}
	testStorage.Close()
}
