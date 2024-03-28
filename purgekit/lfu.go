package purgekit

import "container/list"

type LFUCache struct {
	maxEntries int
	onEnvicted func(key Key, value interface{})

	freqList map[int]*list.List
	cache    map[interface{}]*list.Element
	minFreq  int
}

type lfuEntry struct {
	key   Key
	value interface{}
	freq  int
}

func NewLFUCache(maxEntries int, onEnvicted func(key Key, value interface{})) *LFUCache {
	return &LFUCache{
		maxEntries: maxEntries,
		onEnvicted: onEnvicted,
		freqList:   make(map[int]*list.List),
		cache:      make(map[interface{}]*list.Element),
	}
}

func (c *LFUCache) Get(key Key) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		c.Jump(ele)
		if c.freqList[c.minFreq].Len() == 0 {
			c.minFreq += 1
		}
	}
	return
}

func (c *LFUCache) Add(key Key, value interface{}) {
	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.freqList = make(map[int]*list.List)
	}
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*lfuEntry)
		kv.value = value
		c.Jump(ele)
		if c.freqList[c.minFreq].Len() == 0 {
			c.minFreq += 1
		}
		return
	}
	ll := c.freqList[0]
	if ll == nil {
		ll = list.New()
	}
	ele := ll.PushFront(&lfuEntry{key: key, value: value, freq: 0})
	c.cache[key] = ele
	c.minFreq = 0
	if c.maxEntries != 0 && len(c.cache) > c.maxEntries {
		c.RemoveLeaseUsed()
	}
}

func (c *LFUCache) Remove(key Key) {
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*lfuEntry)
		oldFreq := kv.freq
		c.freqList[oldFreq].Remove(ele)
		delete(c.cache, kv.key)
	}
}

func (c *LFUCache) RemoveLeaseUsed() {
	if c.cache == nil {
		return
	}
	ele := c.freqList[c.minFreq].Back()
	c.freqList[c.minFreq].Remove(ele)
	key := ele.Value.(*lfuEntry).key
	delete(c.cache, key)
}

func (c *LFUCache) Len() int {
	return len(c.cache)
}

// Jump 将一个 ele 提升到频率加一的链表里
func (c *LFUCache) Jump(ele *list.Element) {
	kv := ele.Value.(*lfuEntry)
	oldFreq := kv.freq
	c.freqList[oldFreq].Remove(ele)
	kv.freq += 1
	if ll, ok := c.freqList[kv.freq]; ok {
		ll.PushFront(ele)
		return
	}
	ll := list.New()
	ll.PushFront(ele)
	c.freqList[kv.freq] = ll
}
