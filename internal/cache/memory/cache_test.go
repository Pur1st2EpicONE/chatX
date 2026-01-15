package memory

import (
	"chatX/internal/config"
	"chatX/internal/errs"
	"chatX/internal/logger/mocks"
	"chatX/internal/models"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func testCacheConfig(capacity int, maxMessages int) config.Cache {
	return config.Cache{Capacity: capacity, MaxMessages: maxMessages}
}

func testChat(id int, msgCount int) models.Chat {
	return models.Chat{ID: id, Messages: make([]models.Message, msgCount)}
}

func setupCache(t *testing.T, cap, maxMsgs int) *LRUCache {

	controller := gomock.NewController(t)
	logger := mocks.NewMockLogger(controller)

	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().LogInfo(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().LogWarn(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().LogError(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().LogFatal(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	return NewLRUCache(logger, testCacheConfig(cap, maxMsgs))

}

func TestLRUCache_Get_Disabled(t *testing.T) {
	cache := setupCache(t, 0, 10)
	_, err := cache.Get(1)
	require.ErrorIs(t, err, errs.ErrCacheMiss)
}

func TestLRUCache_Get_Miss(t *testing.T) {
	cache := setupCache(t, 2, 10)
	_, err := cache.Get(42)
	require.ErrorIs(t, err, errs.ErrCacheMiss)
}

func TestLRUCache_PutAndGet_OK(t *testing.T) {
	cache := setupCache(t, 2, 10)
	chat := testChat(1, 1)
	cache.Put(1, chat)
	got, err := cache.Get(1)
	require.NoError(t, err)
	require.Equal(t, chat, got)
}

func TestLRUCache_Put_Disabled(t *testing.T) {
	cache := setupCache(t, 0, 10)
	cache.Put(1, testChat(1, 1))
	_, err := cache.Get(1)
	require.ErrorIs(t, err, errs.ErrCacheMiss)
}

func TestLRUCache_Put_MessageLimitExceeded(t *testing.T) {
	cache := setupCache(t, 2, 1)
	cache.Put(1, testChat(1, 2))
	_, err := cache.Get(1)
	require.ErrorIs(t, err, errs.ErrCacheMiss)
}

func TestLRUCache_Put_Overwrite(t *testing.T) {

	cache := setupCache(t, 2, 10)

	cache.Put(1, testChat(1, 1))
	cache.Put(1, testChat(1, 2))

	chat, err := cache.Get(1)

	require.NoError(t, err)
	require.Len(t, chat.Messages, 2)

}

func TestLRUCache_Put_LRUEviction(t *testing.T) {

	cache := setupCache(t, 2, 10)

	cache.Put(1, testChat(1, 1))
	cache.Put(2, testChat(2, 1))

	_, _ = cache.Get(1)
	cache.Put(3, testChat(3, 1))

	_, err := cache.Get(2)
	require.ErrorIs(t, err, errs.ErrCacheMiss)

	_, err = cache.Get(1)
	require.NoError(t, err)

	_, err = cache.Get(3)
	require.NoError(t, err)

}

func TestLRUCache_Delete_OK(t *testing.T) {

	cache := setupCache(t, 2, 10)

	cache.Put(1, testChat(1, 1))
	cache.Delete(1)

	_, err := cache.Get(1)
	require.ErrorIs(t, err, errs.ErrCacheMiss)

}

func TestLRUCache_Delete_Miss(t *testing.T) {
	cache := setupCache(t, 2, 10)
	cache.Delete(42)
}

func TestLRUCache_Delete_Disabled(t *testing.T) {
	cache := setupCache(t, 0, 10)
	cache.Delete(1)
}

func TestLRUCache_Close(t *testing.T) {

	cache := setupCache(t, 2, 10)

	cache.Put(1, testChat(1, 1))
	cache.Put(2, testChat(2, 1))

	cache.Close()

	require.Nil(t, cache.head)
	require.Nil(t, cache.tail)
	require.Len(t, cache.hm, 0)

}
