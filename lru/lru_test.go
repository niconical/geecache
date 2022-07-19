package lru

import (
	"reflect"
	"testing"
)

type MyString string

func (s MyString) Len() int { return len(s) }

var _ Value = (*MyString)(nil)

func TestLRU_Add_Get(t *testing.T) {
	lru := New(1024, nil)
	lru.Add("key", MyString("value"))
	if v, ok := lru.Get("key"); !ok || string(v.(MyString)) != "value" {
		t.Fatal("cache hit key=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatal("cache miss key2 failed")
	}
}

func TestLRU_RemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	lru := New(int64(len(k1+k2+v1+v2)), nil)
	lru.Add(k1, MyString(v1))
	lru.Add(k2, MyString(v2))
	lru.Add(k3, MyString(v3))
	if _, ok := lru.Get(k1); ok || lru.ll.Len() != 2 {
		t.Fatal("cache remove oldest failed")
	}
}

func TestLRU_OnEvicted(t *testing.T) {
	keys := make([]string, 0)
	lru := New(int64(10), func(k string, v Value) {
		keys = append(keys, k)
	})
	lru.Add("key1", MyString("123456"))
	lru.Add("k2", MyString("k2"))
	lru.Add("k3", MyString("k3"))
	lru.Add("k4", MyString("k4"))
	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
