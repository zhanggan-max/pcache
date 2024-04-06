package purgekit

// ARC 算法使用一个额外的 LRU 缓存保存频率信息
// 根据缓存命中情况,动态调整二者之间的比例
type ARCache struct {
	maxEntries int // maxEntries 保存有效缓存总大小
	onEnvicted func(key Key, value interface{})

	p  int       // p 是缓存中 t1 的长度
	t1 *LRUCache // t1 保存只出现一次的缓存数据
	b1 *LRUCache // b1 保存 t1 中淘汰下来的缓存
	t2 *LRUCache // t2 保存请求超过一次的数据
	b2 *LRUCache // b2 保存 t2 中淘汰的键信息,用来动态调整 p 值
}

// NewARCache 返回一个 ARCache 实例
func NewARCache(maxEntries int, onEnvicted func(key Key, value interface{})) *ARCache {
	return &ARCache{
		maxEntries: maxEntries,
		onEnvicted: onEnvicted,
		p:          0,
		t1:         NewLRUCache(maxEntries),
		t2:         NewLRUCache(maxEntries),
		b1:         NewLRUCache(maxEntries),
		b2:         NewLRUCache(maxEntries),
	}
}

// Get 查找当前缓存,返回对应的缓存值(如果存在),和一个代表操作是否成功的布尔值
func (c *ARCache) Get(key Key) (value interface{}, ok bool) {
	// t1 中找到,移动到 t2
	if val, ok := c.t1.Peek(key); ok {
		c.t1.Remove(key)
		c.t2.Add(key, val)
		return val, true
	}
	// t2 中找到,移动到 t2 队首
	if val, ok := c.t2.Get(key); ok {
		return val, ok
	}
	return
}

// Add 添加一个新缓存或者更新值
func (c *ARCache) Add(key Key, value interface{}) {
	// t1 中找到,移动到 t2
	if c.t1.Contains(key) {
		c.t1.Remove(key)
		c.t2.Add(key, value)
		return
	}
	// t2 中找到,移动到 t2 队首
	if c.t2.Contains(key) {
		c.t2.Add(key, value)
		return
	}

	// b1中找到,说明 t1 太小
	if c.b1.Contains(key) {
		// b1 > b2 时,p 只需要增加 1
		delta := 1
		b1len := c.b1.Len()
		b2len := c.b2.Len()
		// b1 < b2 时, p 需要大幅度调整
		if b2len > b1len {
			delta = b2len / b1len
		}
		// 如果调整后的 t1 长度超过总长度,设置为最大长度
		if c.p+delta >= c.maxEntries {
			c.p = c.maxEntries
		} else {
			c.p += delta
		}
		// 调整后, 如果 t1 + t2 超出总长度, 应当淘汰数据
		if c.t1.Len()+c.t2.Len() >= c.maxEntries {
			c.replace(false)
		}
		// 从淘汰记录中删除当前 key
		c.b1.Remove(key)
		c.t2.Add(key, value)
		return
	}

	// b2 中找到, 说明 t2 过小, 应当减小 p 值
	if c.b2.Contains(key) {
		delta := 1
		b1len := c.b1.Len()
		b2len := c.b2.Len()
		if b1len > b2len {
			delta = b1len / b2len
		}
		// 调整后 p 小于 0, 设置为 0
		if delta >= c.p {
			c.p = 0
		} else {
			c.p -= delta
		}
		// 调整后有效缓存超过总长度, 淘汰数据
		if c.t1.Len()+c.t2.Len() >= c.maxEntries {
			c.replace(true)
		}
		c.b2.Remove(key)
		c.t2.Add(key, value)
		return
	}
	// t1 t2 b1 b2 都没有, 需要添加新条目
	// 如果没有空间, 先淘汰
	if c.t1.Len()+c.t2.Len() >= c.maxEntries {
		c.replace(false)
	}
	if c.b1.Len() > c.maxEntries-c.p {
		c.b1.RemoveOldest()
	}
	if c.b2.Len() > c.p {
		c.b2.RemoveOldest()
	}
	// 添加新条目
	c.t1.Add(key, value)
}

// replace 根据情况选择不同队列淘汰
// 如果 t1 超出当前 p 值, 并且 t2 过小,从 t1 中淘汰
// 否则从 t2 中淘汰
func (c *ARCache) replace(contains bool) {
	t1len := c.t1.Len()
	if t1len > 0 && (t1len > c.p || (t1len == c.p && contains)) {
		k, _, ok := c.t1.RemoveOldest()
		if ok {
			c.b1.Add(k, nil)
		}
	} else {
		k, _, ok := c.t2.RemoveOldest()
		if ok {
			c.b2.Add(k, nil)
		}
	}
}

// Len 返回当期有效缓存的数量
func (c *ARCache) Len() int {
	return c.t1.Len() + c.t2.Len()
}

func (c *ARCache) RegisterOnEnvicted(onfunc OnEnvictedFunc) {
	c.onEnvicted = onfunc
}
