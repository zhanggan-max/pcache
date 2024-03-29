package pcache

import (
	"pcache/purgekit"
	"sync"
)

type cache struct {
	m          sync.RWMutex
	lru        purgekit.Cache
	maxEntries int
}

func (c *cache) add(key string, value ByteView) {
	c.m.RLock()
	defer c.m.RUnlock()
	if c.lru == nil {
		panic("please init cache first")
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.m.RLock()
	defer c.m.RUnlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
