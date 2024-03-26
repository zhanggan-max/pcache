package purgekit

import (
	"fmt"
	"testing"
)

type simpleStruct struct {
	string
	int
}

type complextStruct struct {
	string
	simpleStruct
}

var getTests = []struct {
	name     string
	keyToAdd interface{}
	keyToGet interface{}
	expected bool
}{
	{"string_hit", "myKey", "myKey", true},
	{"string_miss", "myKey", "nonsnese", false},
	{"simpleStruct_hit", simpleStruct{"mykey", 1}, simpleStruct{"mykey", 1}, true},
	{"simpleStruct_miss", simpleStruct{"mykey", 1}, simpleStruct{"mykey", 2}, false},
	{"comlextStruct_hit", complextStruct{"key", simpleStruct{"key", 1}}, complextStruct{"key", simpleStruct{"key", 1}}, true},
}

func TestGet(t *testing.T) {
	for _, test := range getTests {
		lru := NewLRUCache(0)
		lru.Add(test.keyToAdd, 1234)
		value, ok := lru.Get(test.keyToGet)
		if ok != test.expected {
			t.Fatalf("%s: cache hit %v; want %v", test.name, ok, !ok)
		} else if ok && value != 1234 {
			t.Fatalf("%s want 1234, but got %v", test.name, value)
		}
	}
}

func TestRemove(t *testing.T) {
	lru := NewLRUCache(0)
	lru.Add("key", 1234)
	if value, ok := lru.Get("key"); !ok {
		t.Fatal("TestRemove returns no match!")
	} else if value != 1234 {
		t.Fatalf("TestRemove expected %v but got %v", 1234, value)
	}

	lru.Remove("key")
	if _, ok := lru.Get("key"); ok {
		t.Fatal("TestRemove returned a removed entry!")
	}
}

func TestOnEvictd(t *testing.T) {
	envictedKeys := make([]Key, 0)
	envictedFunc := func(key Key, value interface{}) {
		envictedKeys = append(envictedKeys, key)
	}

	lru := NewLRUCache(20)
	lru.onEnvited = envictedFunc
	for i := 0; i < 22; i++ {
		lru.Add(fmt.Sprintf("mykey%d", i), 1234)
	}
	if len(envictedKeys) != 2 {
		t.Fatalf("got %d envicted keys, want 22", len(envictedKeys))
	}
	if envictedKeys[0] != Key("mykey0") {
		t.Fatalf("got %v in first envicted key; want mykey0", envictedKeys[0])
	}
	if envictedKeys[1] != Key("mykey1") {
		t.Fatalf("got %v in first envicted key; want mykey0", envictedKeys[1])
	}
}
