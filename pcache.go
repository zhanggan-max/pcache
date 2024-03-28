package pcache

import (
	"fmt"
	"log"
	"pcache/singleflight"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
	server    Picker
	flight    *singleflight.Flight
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, maxEntries int, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
	}
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{maxEntries: maxEntries},
		flight:    &singleflight.Flight{},
	}
	mu.Lock()
	groups[name] = g
	mu.Unlock()
	return g
}

func (g *Group) RegisterPicker(p Picker) {
	if g.server != nil {
		panic("group had been registered server")
	}
	g.server = p
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

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

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populate(key, value)
	return value, nil
}

func (g *Group) populate(key string, value ByteView) {
	g.mainCache.add(key, value)
}
