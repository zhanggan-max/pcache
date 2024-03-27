package purgekit

import (
	"container/list"
)

type ARCache struct {
	maxEntries int
	onEnvited  func(key Key, value interface{})

	t1 *list.List
	b1 *list.List

	t2 *list.List
	b2 *list.List

	cache map[interface{}]*list.Element
}

type arcEntry struct {
	key    Key
	value  interface{}
	parent *list.List
}

func newARCache(maxEntries int, onEnvited func(key Key, value interface{})) *ARCache {
	return &ARCache{
		maxEntries: maxEntries,
		onEnvited:  onEnvited,
		t1:         list.New(),
		t2:         list.New(),
		b1:         list.New(),
		b2:         list.New(),
		cache:      make(map[interface{}]*list.Element),
	}
}

func (c *ARCache) Get(key Key) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*arcEntry)
		switch kv.parent {
		case c.t1:
			c.t1.Remove(ele)
			c.t2.PushFront(ele)
			kv.parent = c.t2
		case c.t2:
			c.t2.MoveToFront(ele)
		case c.b1:
			c.b1.Remove(ele)
			c.t2.PushFront(ele)
			kv.parent = c.t2
		case c.b2:
			c.b2.Remove(ele)
			c.t2.PushFront(ele)
			kv.parent = c.t2
		}
		if c.t1.Len()+c.t2.Len() > c.maxEntries/2 {
			if c.b1.Len()+c.b2.Len() > c.maxEntries/2 {
				if c.b1 == nil || c.b1.Len() == 0 {
					back := c.b2.Back()
					c.b2.Remove(back)
					back2 := c.t2.Back()
					c.t2.Remove(back2)
					c.b2.PushFront(back2)
				} else {
					back := c.b1.Back()
					c.b1.Remove(back)
					c.b2.PushFront(c.t2.Remove(c.t2.Back()))
				}
			} else {
				c.b2.PushFront(c.t2.Remove(c.t2.Back()))
			}
		}
		return kv.value, true
	}
	return
}

// 回调函数没有使用
func (c *ARCache) Add(key Key, value interface{}) {
	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.t1 = list.New()
		c.t2 = list.New()
		c.b1 = list.New()
		c.b2 = list.New()
	}
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*arcEntry)
		kv.value = value
		c.Get(key)
	}
	ele := c.t1.PushFront(&arcEntry{key: key, value: value, parent: c.t1})
	c.cache[key] = ele
	if c.t1.Len()+c.t2.Len() > c.maxEntries/2 {
		if c.b1.Len()+c.b2.Len() < c.maxEntries/2 {
			c.b1.PushFront(c.t1.Remove(c.t1.Back()))
		} else {
			if c.b1 == nil || c.b1.Len() == 0 {
				c.b2.Remove(c.b2.Back())
				c.b1.PushFront(c.t1.Remove(c.t1.Back()))
			} else {
				c.b1.Remove(c.b1.Back())
				c.b1.PushFront(c.t1.Remove(c.t1.Back()))
			}
		}
	}
}

func (c *ARCache) Remove(key Key) {
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*arcEntry)
		delete(c.cache, key)
		kv.parent.Remove(ele)
	}
}

func (c *ARCache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	if c.b1 == nil || c.b1.Len() == 0 {
		c.b2.Remove(c.b2.Back())
		return
	}
	c.b1.Remove(c.b1.Back())
}

func (c *ARCache) Len() int {
	return c.t1.Len() + c.t2.Len()
}

func (c *ARCache) Total() int {
	return c.Len() + c.b1.Len() + c.b2.Len()
}

func (c *ARCache) Clear() {
	c.cache = nil
	c.t1 = nil
	c.t2 = nil
	c.b1 = nil
	c.b2 = nil
}
