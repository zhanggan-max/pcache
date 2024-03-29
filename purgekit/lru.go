package purgekit

import "container/list"

// LRUCache 是实现了 LRU 淘汰机制的结构
type LRUCache struct {
	maxEntries int                              // maxEntries 是缓存可存储的最大条目数
	onEnvited  func(key Key, value interface{}) // 淘汰条目时进行的额外操作

	ll    *list.List
	cache map[interface{}]*list.Element
}

// entry 是在 list 中使用的结构
type entry struct {
	key   Key
	value interface{}
}

func NewLRUCache(maxEntries int) *LRUCache {
	return &LRUCache{
		maxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[interface{}]*list.Element),
	}
}

// Get 返回缓存中对应的值（如果存在）和一个表示是否存在的布尔值
func (c *LRUCache) Get(key Key) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return
}

// Add 向缓存中添加键值对（不存在）或更新键值对（已存在）
func (c *LRUCache) Add(key Key, value interface{}) {
	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.ll = list.New()
	}
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		ele.Value.(*entry).value = value
		return
	}
	ele := c.ll.PushFront(&entry{key, value})
	c.cache[key] = ele
	if c.maxEntries != 0 && c.ll.Len() > c.maxEntries {
		c.RemoveOldest()
	}
}

// Remove 移除给定键对应的缓存
func (c *LRUCache) Remove(key Key) {
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		c.removeElement(ele)
	}
}

// 从缓存中移除缓存中最老的条目
func (c *LRUCache) RemoveOldest() (key Key, value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	kv := ele.Value.(*entry)
	if ele != nil {
		c.removeElement(ele)
	}
	return kv.key, kv.value, true
}

// removeElement 从缓存中删除指定的缓存
func (c *LRUCache) removeElement(ele *list.Element) {
	c.ll.Remove(ele)
	kv := ele.Value.(*entry)
	delete(c.cache, kv.key)
	if c.onEnvited != nil {
		c.onEnvited(kv.key, kv.value)
	}
}

// 返回当前缓存中的条目数
func (c *LRUCache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

// Clear 移除所有的缓存
func (c *LRUCache) Clear() {
	if c.cache == nil {
		return
	}
	c.cache = nil
	c.ll = nil
}

func (c *LRUCache) Contains(key Key) bool {
	if c.cache == nil {
		return false
	}
	_, ok := c.cache[key]
	return ok
}

func (c *LRUCache) Peek(key Key) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		return ele.Value.(*entry).value, true
	}
	return
}
