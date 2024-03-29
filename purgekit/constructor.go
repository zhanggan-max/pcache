package purgekit

// 运训任意可比较的类型作为键
type Key interface{}

// OnEnvictedFunc 是可选的参数，在移除缓存元素是调用
type OnEnvictedFunc func(Key, interface{})

// Cache 接口向外开放
type Cache interface {
	Get(Key) (interface{}, bool)
	Add(Key, interface{})
	Len() int
}

// NewCache 根据 policy 选择实例化对应的缓存
// 错误的 policy 将会返回 LRUCache 实例
func NewCache(policy string, maxEntries int, onEnvicted OnEnvictedFunc) Cache {
	switch policy {
	case "lru":
		return &LRUCache{maxEntries: maxEntries, onEnvited: onEnvicted}
	case "lfu":
		return &LFUCache{maxEntries: maxEntries, onEnvicted: onEnvicted}
	case "arc":
		return &ARCache{maxEntries: maxEntries, onEnvicted: onEnvicted}
	default:
		return &LRUCache{maxEntries: maxEntries, onEnvited: onEnvicted}
	}
}
