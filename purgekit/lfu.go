package purgekit

import "container/list"

type LFUCache struct {
	maxEntries int
	onEnvited  func(key Key, value interface{})

	// freqList 存储使用频率和条目的对应关系
	// 使用双向链表存储全部条目，降低增加和删除的成本
	freqList map[int]*list.List
	// cache 存储键值对应，降低查询的成本
	cache map[interface{}]*list.Element
	// minFrequency 代表一个一个缓存第一次添加时的 frequency
	minFrequency int
}

// 在节点中存储的条目
type lfuEntry struct {
	key       Key
	value     interface{}
	frequency int // freq 代表当前条目被查询的次数
}

func NewLFUCache(maxEntries int, onEnvited func(key Key, value interface{})) *LFUCache {
	return &LFUCache{
		maxEntries:   maxEntries,
		onEnvited:    onEnvited,
		minFrequency: 0,
		freqList:     make(map[int]*list.List),
		cache:        make(map[interface{}]*list.Element),
	}
}

func (c *LFUCache) Get(key Key) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*lfuEntry)
		c.freqList[kv.frequency].Remove(ele)
		kv.frequency += 1
		if ll, ok := c.freqList[kv.frequency]; ok {
			ll.PushFront(ele)
		} else {
			ll := list.New()
			ll.PushFront(ele)
			c.freqList[kv.frequency] = ll
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
		c.freqList[kv.frequency].Remove(ele)
		kv.frequency += 1
		if ll, ok := c.freqList[kv.frequency]; ok {
			ll.PushFront(ele)
		} else {
			ll := list.New()
			ll.PushFront(ele)
			c.freqList[kv.frequency] = ll
		}
	}
	entry := &lfuEntry{key: key, value: value, frequency: c.minFrequency}
	if ll, ok := c.freqList[entry.frequency]; ok {
		ele := ll.PushFront(entry)
		c.cache[key] = ele
	}
	ll := list.New()
	ele := ll.PushFront(entry)
	c.cache[key] = ele
	if c.maxEntries != 0 && c.Len() > c.maxEntries {
		c.RemoveOldest()
	}
}

func (c *LFUCache) Remove(key Key) {
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*lfuEntry)
		delete(c.cache, key)
		c.freqList[kv.frequency].Remove(ele)
	}
}

func (c *LFUCache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	ll := c.freqList[c.minFrequency]
	ele := ll.Back()
	ll.Remove(ele)
	if c.onEnvited != nil {
		kv := ele.Value.(*lfuEntry)
		c.onEnvited(kv.key, kv.value)
	}
}

func (c *LFUCache) Len() int {
	return len(c.cache)
}

func (c *LFUCache) Clear() {
	c.cache = nil
	c.freqList = nil
}
