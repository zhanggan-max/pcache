package purgekit

type ARCache struct {
	maxEntries int
	onEnvicted func(key Key, value interface{})

	p  int
	t1 *LRUCache
	b1 *LRUCache
	t2 *LFUCache
	b2 *LFUCache
}

func NewARCache(maxEntries int, onEnvicted func(key Key, value interface{})) *ARCache {
	return &ARCache{
		maxEntries: maxEntries,
		onEnvicted: onEnvicted,
		p:          0,
		t1:         NewLRUCache(maxEntries),
		t2:         NewLFUCache(maxEntries, onEnvicted),
		b1:         NewLRUCache(maxEntries),
		b2:         NewLFUCache(maxEntries, onEnvicted),
	}
}

func (c *ARCache) Get(key Key) (value interface{}, ok bool) {
	if val, ok := c.t1.Peek(key); ok {
		c.t1.Remove(key)
		c.t2.Add(key, val)
		return val, true
	}
	if val, ok := c.t2.Get(key); ok {
		return val, ok
	}
	return
}

func (c *ARCache) Add(key Key, value interface{}) {
	if c.t1.Contains(key) {
		c.t1.Remove(key)
		c.t2.Add(key, value)
		return
	}
	if c.t2.Contains(key) {
		c.t2.Add(key, value)
		return
	}

	if c.b1.Contains(key) {
		delta := 1
		b1len := c.b1.Len()
		b2len := c.b2.Len()
		if b2len > b1len {
			delta = b2len / b1len
		}
		if c.p+delta >= c.maxEntries {
			c.p = c.maxEntries
		} else {
			c.p += delta
		}
		if c.t1.Len()+c.t2.Len() >= c.maxEntries {
			c.replace(false)
		}
		c.b1.Remove(key)
		c.t2.Add(key, value)
		return
	}

	if c.b2.Contains(key) {
		delta := 1
		b1len := c.b1.Len()
		b2len := c.b2.Len()
		if b1len > b2len {
			delta = b1len / b2len
		}
		if delta >= c.p {
			c.p = 0
		} else {
			c.p -= delta
		}
		if c.t1.Len()+c.t2.Len() >= c.maxEntries {
			c.replace(true)
		}
		c.b2.Remove(key)
		c.t2.Add(key, value)
		return
	}
	if c.t1.Len()+c.t2.Len() >= c.maxEntries {
		c.replace(false)
	}
	if c.b1.Len() > c.maxEntries-c.p {
		c.b1.RemoveOldest()
	}
	if c.b2.Len() > c.p {
		c.b2.RemoveLeastUsed()
	}
	c.t1.Add(key, value)
}

func (c *ARCache) replace(contains bool) {
	t1len := c.t1.Len()
	if t1len > 0 && (t1len > c.p || contains) {
		k, _, ok := c.t1.RemoveOldest()
		if ok {
			c.b1.Add(k, nil)
		}
	} else {
		k, _, ok := c.t2.RemoveLeastUsed()
		if ok {
			c.b2.Add(k, nil)
		}
	}
}

func (c *ARCache) Len() int {
	return c.t1.Len() + c.t2.Len()
}
