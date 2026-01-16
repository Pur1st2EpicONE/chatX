// Package memory provides in-memory cache implementations,
// including an LRU cache for storing chats.
package memory

import (
	"chatX/internal/config"
	"chatX/internal/errs"
	"chatX/internal/logger"
	"chatX/internal/models"
	"sync"
)

// Node represents a doubly-linked list node for LRUCache.
type Node struct {
	Key  int         // Cache key (chat ID)
	Val  models.Chat // Cached chat
	Next *Node       // Pointer to the next node
	Prev *Node       // Pointer to the previous node
}

// newNode creates a new linked-list node for the given key and chat.
func newNode(key int, value models.Chat) *Node {
	return &Node{Key: key, Val: value}
}

// LRUCache is a thread-safe in-memory LRU cache for chats.
type LRUCache struct {
	mu     sync.RWMutex  // Mutex for concurrent access
	head   *Node         // Dummy head node
	tail   *Node         // Dummy tail node
	hm     map[int]*Node // Map of keys to nodes
	config config.Cache  // Cache configuration
	logger logger.Logger // Logger instance
}

// NewLRUCache creates a new LRUCache instance with the given logger and config.
func NewLRUCache(logger logger.Logger, config config.Cache) *LRUCache {
	head := newNode(0, models.Chat{})
	tail := newNode(0, models.Chat{})
	head.Next = tail
	tail.Prev = head
	return &LRUCache{
		head:   head,
		tail:   tail,
		hm:     make(map[int]*Node, config.Capacity),
		config: config,
		logger: logger,
	}
}

// remove deletes a node from the linked list and map.
func (c *LRUCache) remove(node *Node) {
	delete(c.hm, node.Key)
	node.Next.Prev = node.Prev
	node.Prev.Next = node.Next
	node.Prev, node.Next = nil, nil
}

// insert adds a node to the front of the linked list and updates the map.
func (c *LRUCache) insert(node *Node) {
	c.hm[node.Key] = node
	next := c.head.Next
	c.head.Next = node
	node.Prev = c.head
	node.Next = next
	next.Prev = node
}

// Get retrieves a chat from the cache by key and moves it to the front (most recently used).
func (c *LRUCache) Get(key int) (models.Chat, error) {

	if c.config.Capacity <= 0 {
		return models.Chat{}, errs.ErrCacheMiss
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if node, ok := c.hm[key]; ok {
		c.logger.Debug("cache — chat found", "chatID", key, "layer", "cache.memory")
		c.remove(node)
		c.insert(node)
		return node.Val, nil
	}
	c.logger.Debug("cache — chat not found", "chatID", key, "layer", "cache.memory")

	return models.Chat{}, errs.ErrCacheMiss

}

// Put stores a chat in the cache. Evicts least-recently-used chat if capacity is exceeded.
func (c *LRUCache) Put(key int, value models.Chat) {

	if c.config.Capacity <= 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if len(value.Messages) > c.config.MaxMessages {
		c.logger.Debug("cache — chat not cached: message limit exceeded", "layer", "cache.memory")
		return
	}

	if node, ok := c.hm[key]; ok {
		c.remove(node)
	}

	if len(c.hm) == c.config.Capacity {
		c.logger.Debug("cache — maximum capacity reached", "layer", "cache.memory")
		lru := c.tail.Prev
		c.remove(lru)
		c.logger.Debug("cache — LRU chat deleted", "chatID", lru.Key, "layer", "cache.memory")
	}

	c.insert(newNode(key, value))
	c.logger.Debug("cache — chat saved", "chatID", key, "layer", "cache.memory")

}

// Delete removes a chat from the cache by key.
func (c *LRUCache) Delete(key int) {

	if c.config.Capacity <= 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if node, ok := c.hm[key]; ok {
		c.remove(node)
		c.logger.Debug("cache — chat deleted", "chatID", key, "layer", "cache.memory")
	}

}

// Close releases all resources used by the cache.
func (c *LRUCache) Close() {

	for k := range c.hm {
		delete(c.hm, k)
	}

	c.head = nil
	c.tail = nil

	c.logger.LogInfo("cache — resources released", "layer", "cache.memory")

}
