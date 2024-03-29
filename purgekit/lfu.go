package purgekit

import "container/list"

type LFUCache struct {
	maxEntries int
	onEnvicted func(key Key, value interface{})

	freqList map[int]*list.List
	cache    map[interface{}]*list.Element
	minFreq  int
}

// lfuEntry 存储具体的条目，在链表节点中使用
type lfuEntry struct {
	key   Key
	value interface{}
	freq  int
}

// NewLFUCache 返回一个 lfucache 对象指针
func NewLFUCache(maxEntries int, onEnvicted func(key Key, value interface{})) *LFUCache {
	return &LFUCache{
		maxEntries: maxEntries,
		onEnvicted: onEnvicted,
		freqList:   make(map[int]*list.List),
		cache:      make(map[interface{}]*list.Element),
		minFreq:    0,
	}
}

// Get 返回 key 对应的值（如果存在），和一个表示值是否存在的布尔值
func (c *LFUCache) Get(key Key) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		c.Jump(ele)
		kv := ele.Value.(*lfuEntry)
		// todo: 这里有一个 bug
		if kv.freq == c.minFreq && c.freqList[c.minFreq].Len() == 0 {
			c.minFreq += 1
		}
		return ele.Value.(*lfuEntry).value, true
	}
	return
}

// Add 添加一个键值对到缓存中（如果之前不存在），或者更新一个键值对（之前已经存在）
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
	c.minFreq = 0
	ll := c.freqList[0]
	if ll == nil {
		ll = list.New()
		c.freqList[c.minFreq] = ll
	}
	ele := ll.PushFront(&lfuEntry{key: key, value: value, freq: 0})
	c.cache[key] = ele
	if c.maxEntries != 0 && len(c.cache) > c.maxEntries {
		c.RemoveLeastUsed()
	}
}

// Remove 移除指定的 Key
// todo: 如果 remove 移除的是将要淘汰的最后一个值， minfreq 如何变化
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

// RemoveLeastUsed 移除使用频率最少，使用时间最久远的键值对
func (c *LFUCache) RemoveLeastUsed() (key Key, value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	// todo: 这里有一个 bug
	ele := c.freqList[c.minFreq].Back()
	c.freqList[c.minFreq].Remove(ele)
	kv := ele.Value.(*lfuEntry)
	delete(c.cache, kv.key)
	return kv.key, kv.value, true
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

func (c *LFUCache) Contains(key Key) bool {
	if c.cache == nil {
		return false
	}
	if _, ok := c.cache[key]; ok {
		return true
	}
	return false
}

func (c *LFUCache) Peek(key Key) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		return ele.Value.(*lfuEntry).value, true
	}
	return
}
