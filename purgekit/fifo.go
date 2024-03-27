package purgekit

import "container/list"

type FIFOCache struct {
	maxEntries int
	onEnvited  func(key Key, value interface{})

	ll    *list.List
	cache map[interface{}]*list.Element
}

type fifoEntry struct {
	key   Key
	value interface{}
}

func newFIFOCache(maxEntries int, onEnvited func(key Key, value interface{})) *FIFOCache {
	return &FIFOCache{
		maxEntries: maxEntries,
		onEnvited:  onEnvited,
		ll:         list.New(),
		cache:      make(map[interface{}]*list.Element),
	}
}

func (c *FIFOCache) Get(key Key) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		return ele.Value.(*fifoEntry).value, true
	}
	return
}

func (c *FIFOCache) Add(key Key, value interface{}) {
	if c.cache == nil {
		c.ll = list.New()
		c.cache = make(map[interface{}]*list.Element)
	}
	ele := c.ll.PushFront(&fifoEntry{key: key, value: value})
	c.cache[key] = ele
	if c.maxEntries != 0 && c.Len() > c.maxEntries {
		key, value := c.RemoveFirstIn()
		c.onEnvited(key, value)
	}
}

func (c *FIFOCache) Remove(key Key) {
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		delete(c.cache, key)
		kv := c.ll.Remove(ele).(*fifoEntry)
		c.onEnvited(kv.key, kv.value)
	}
}

func (c *FIFOCache) RemoveFirstIn() (key interface{}, value interface{}) {
	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	c.ll.Remove(ele)
	delete(c.cache, ele)
	kv := ele.Value.(*fifoEntry)
	return kv.key, kv.value
}

func (c *FIFOCache) Len() int {
	return c.ll.Len()
}

func (c *FIFOCache) Clear() {
	c.cache = nil
	c.ll = nil
}
