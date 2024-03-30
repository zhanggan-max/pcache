package pcache

import (
	"fmt"
	"log"
	"pcache/singleflight"
	"sync"
)

// Getter 接口包含一个从数据源获取数据的 Get 方法
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 实现了 Getter 接口，方便编写获取函数
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group 提供了用户的交互入口
type Group struct {
	name      string               // name 是当前 Group 的名字
	getter    Getter               // getter 从数据源获得数据
	mainCache *cache               // mainCache 是真正的缓存
	server    Picker               // server 从注册节点中选择节点
	flight    *singleflight.Flight // flight 确保一个键同时只有一次请求
}

var (
	mu sync.RWMutex
	// groups 管理当前所有的 Group，是并发安全的
	groups = make(map[string]*Group)
)

// NewGroup 创建一个 Group 实例，并注册到 groups 中
func NewGroup(name string, maxEntries int, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
	}
	g := &Group{
		name:   name,
		getter: getter,
		// mainCache 的初始化未完成
		mainCache: &cache{maxEntries: maxEntries},
		flight:    &singleflight.Flight{},
	}
	mu.Lock()
	groups[name] = g
	mu.Unlock()
	return g
}

// RegisterPicker 将节点选择器注册到 Group
func (g *Group) RegisterPicker(p Picker) {
	if g.server != nil {
		panic("group had been registered server")
	}
	g.server = p
}

// GetGroup 获取 name 对应的 Group
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// DestroyGroup 将 name 对饮的 Group 下线
func DestroyGroup(name string) {
	g := GetGroup(name)
	if g != nil {
		picker := g.server.(*server)
		picker.Stop()
		mu.Lock()
		delete(groups, name)
		mu.Unlock()
		log.Printf("Destroy cache %s %s", name, picker.addr)
	}
}

// Get 尝试从当前节点获取 key 对应的值
// 如果本地不存在，尝试从其他节点获得
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.get(key); ok {
		log.Println("Pcache hit")
		return v, nil
	}
	return g.load(key)
}

// load 使用 flight 保证同一个 key 不会多次请求
// 如果远程节点当前也没有缓存，会调用 getter 从数据源获取
func (g *Group) load(key string) (value ByteView, err error) {
	view, err := g.flight.Fly(key, func() (interface{}, error) {
		if g.server != nil {
			if fetcher, ok := g.server.Pick(key); ok {
				bytes, err := fetcher.Fetch(g.name, key)
				if err == nil {
					return ByteView{b: cloneBytes(bytes)}, nil
				}
				log.Printf("failed to get %s from peer, %s\n", key, err.Error())
			}
		}
		return g.getLocally(key)
	})
	if err == nil {
		return view.(ByteView), err
	}
	return ByteView{}, nil
}

// getLocally 从数据源获取数据
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populate(key, value)
	return value, nil
}

// populate 像缓存中添加数据
func (g *Group) populate(key string, value ByteView) {
	g.mainCache.add(key, value)
}

// Name 返回当前 Group 的名字
func (g *Group) Name() string {
	return g.name
}
