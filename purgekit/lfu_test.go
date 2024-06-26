package purgekit

import (
	"testing"
)

func TestLFUAdd(t *testing.T) {
	lfu := NewLFUCache(0, nil)
	lfu.Add("mykey", "value")
	if lfu.Len() != 1 {
		t.Fatalf("Add one key but got %v", lfu.Len())
	}
}

func TestLFUGet(t *testing.T) {
	lfu := NewLFUCache(0, nil)
	lfu.Add("mykey", 1234)
	value, ok := lfu.Get("mykey")
	if !ok || value != 1234 {
		t.Fatalf("key mykey want 1234, but got %v", value)
	}
}

func TestRemoveLeastUsed(t *testing.T) {
	lfu := NewLFUCache(1, nil)
	lfu.Add("key1", 12345)
	lfu.Add("key2", 34)
	if lfu.Len() != 1 {
		t.Fatalf("lfu should have 1 key kept, but got %v", lfu.Len())
	}
	if value, ok := lfu.Get("key1"); ok {
		t.Fatalf("key1 should have been removed, but got reply: %v", value)
	}
}

func TestFreq(t *testing.T) {
	lfu := NewLFUCache(0, nil)
	lfu.Add("key1", 1)
	lfu.Add("key2", 2)
	if lfu.minFreq != 0 {
		t.Fatalf("minFreq should be 0, but got: %v", lfu.minFreq)
	}

	lfu.Get("key1")
	if lfu.minFreq != 0 {
		t.Fatalf("minFreq should be 0, but gout %v", lfu.minFreq)
	}

	lfu.Get("key2")
	if lfu.minFreq != 1 {
		t.Fatalf("minFreq should be 1, but got: %v", lfu.minFreq)
	}
}
