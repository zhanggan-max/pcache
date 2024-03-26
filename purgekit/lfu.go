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
}

// 在节点中存储的条目
type lfuEntry struct {
	key   Key
	value interface{}
	freq  int
}

func newLFUCache(maxEntries int, onEnvited func(key Key, value interface{})) *LFUCache {
	return &LFUCache{
		maxEntries: maxEntries,
		onEnvited:  onEnvited,
		freqList:   make(map[int]*list.List),
		cache:      make(map[interface{}]*list.Element),
	}
}

func (c *LFUCache) Get(key Key) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		c.increaseFreq(ele)
	}
	return
}

func (c *LFUCache) increaseFreq(ele *list.Element) {
	kv := ele.Value.(*lfuEntry)
	freq := kv.freq
	oldlist := c.freqList[freq]
	oldlist.Remove(ele)
	freq += 1
	if _, ok := c.freqList[freq]; !ok {
		c.freqList[freq] = list.New()
	}
	newlist := c.freqList[freq]
	kv.freq = freq
	newlist.PushFront(kv)
}

func (c *LFUCache) Add(key Key, value interface{}) {
	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.freqList = make(map[int]*list.List)
	}
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*lfuEntry)
		kv.value = value
		c.increaseFreq(ele)
		return
	}
	kv := &lfuEntry{key: key, value: value, freq: 0}
	ele := c.freqList[0].PushFront(kv)
	c.cache[key] = ele
}

func (c *LFUCache) Remove(key Key) {
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*lfuEntry)
		delete(c.cache, kv.key)
		c.freqList[kv.freq].Remove(ele)
	}
}

func (c *LFUCache) RemoveLeastFreq() {
	if c.cache == nil {
		return
	}
}
