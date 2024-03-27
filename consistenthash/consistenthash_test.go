package consistenthash

import (
	"hash/crc32"
	"log"
	"sort"
	"testing"
)

func TestRegister(t *testing.T) {
	c := New(2, nil)
	c.Registe("peer1", "peer2")
	if len(c.ring) != 4 {
		t.Errorf("got %d; expect %d", len(c.ring), 4)
	}
	hashValue := int(crc32.ChecksumIEEE([]byte("1peer1")))
	idx := sort.SearchInts(c.ring, hashValue)
	if c.ring[idx] != hashValue {
		t.Errorf("got %d; expect %d", c.ring[idx], hashValue)
	}
}

func TestGetPeer(t *testing.T) {
	c := New(1, nil)
	c.Registe("peer1", "peer2")
	key := "TOM"
	keyHashValue := int(crc32.ChecksumIEEE([]byte(key)))
	log.Printf("key hash = %d", keyHashValue)
	for _, v := range c.ring {
		log.Printf("%d -> %s\n", v, c.hashMap[v])
	}
	peer := c.GetPeer(key)
	log.Printf("Go to search -> %s\n", peer)
}
