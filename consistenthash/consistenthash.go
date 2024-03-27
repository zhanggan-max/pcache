package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// HashFunc 定义哈希函数的输入输出
type HashFunc func(data []byte) uint32

type Map struct {
	hash     HashFunc
	replicas int            // 虚拟节点倍数
	ring     []int          // uint32 哈希环
	hashMap  map[int]string // hashvalue 到 节点之间的映射
}

func New(replicas int, fn HashFunc) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Registe 将节点注册到哈希环上
func (m *Map) Registe(peers ...string) {
	for _, peerName := range peers {
		for i := 0; i < m.replicas; i++ {
			hashValue := int(m.hash([]byte(strconv.Itoa(i) + peerName)))
			m.ring = append(m.ring, hashValue)
			m.hashMap[hashValue] = peerName
		}
	}
	sort.Ints(m.ring)
}

// GetPeer 根据 key 计算应当使用的节点
func (m *Map) GetPeer(key string) string {
	if len(m.ring) == 0 {
		return ""
	}
	hashValue := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.ring), func(i int) bool {
		return m.ring[i] >= hashValue
	})
	return m.hashMap[m.ring[idx%len(m.ring)]]
}
