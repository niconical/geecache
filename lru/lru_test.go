package lru

import (
	"reflect"
	"testing"
)

var _ Value = (*testValueImpl)(nil)

type testValueImpl struct{}

func (v testValueImpl) Len() int64 { return 2 }

var (
	testk1, testk2, testk3 string        = "k1", "k2", "k3"
	testv1, testv2, testv3 testValueImpl = testValueImpl{}, testValueImpl{}, testValueImpl{}
)

// TODO: benchmark

func TestLRU(t *testing.T) {
	OnEvictedCounter := 0
	OnEvicted := func(_ string, _ Value) {
		OnEvictedCounter++
	}

	cache, err := NewCache(12, OnEvicted)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	cache.Add(testk1, testv1)
	cache.Add(testk2, testv2)
	cache.Add(testk3, testv3)
	if cache.Len() != 12 {
		t.Fatalf("bad len: %v", cache.Len())
	}
	if OnEvictedCounter != 0 {
		t.Fatalf("bad evict count: %v", OnEvictedCounter)
	}

	if !reflect.DeepEqual([]string{testk1, testk2, testk3}, cache.Keys()) {
		t.Fatalf("bad keys")
	}

	if ok := cache.Remove(testk1); !ok {
		t.Fatalf("bad remove k1")
	}

	cache.Purge()
	if _, _, ok := cache.GetOldest(); ok {
		t.Fatalf("should get nothing")
	}
	var testKeys = []string{testk1, testk2, testk3}
	for _, k := range testKeys {
		if _, ok := cache.Get(k); ok {
			t.Fatalf("should get nothing")
		}
		if _, ok := cache.Peek(k); ok {
			t.Fatalf("should peek nothing")
		}
	}
	if cache.Len() != 0 {
		t.Fatalf("bad len: %v", cache.Len())
	}
}

func TestLRU_Peek(t *testing.T) {
	cache, err := NewCache(8, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	cache.Add(testk1, testv1)
	cache.Add(testk2, testv2)
	if v, ok := cache.Peek(testk1); !ok || v != testv1 {
		t.Errorf("1 should be set to 1: %v, %v", v, ok)
	}

	cache.Add(testk3, testv3)
	if cache.Contains(testk1) {
		t.Errorf("should not have updated recent-ness of 1")
	}
}

func TestLRU_Contain(t *testing.T) {
	cache, _ := NewCache(8, nil)
	cache.Add(testk1, testv1)
	cache.Add(testk2, testv2)
	if !cache.Contains(testk1) {
		t.Errorf("k1 should be containted")
	}
	cache.Add(testk3, testv3)
	if cache.Contains(testk1) {
		t.Errorf("Contains should not have updated recent-ness of k1")
	}
}

func TestLRU_GetOldest(t *testing.T) {
	cache, _ := NewCache(8, nil)
	cache.Add(testk1, testv1)
	cache.Add(testk2, testv2)
	k, v, ok := cache.GetOldest()
	if k != testk1 || v != testv1 || !ok {
		t.Errorf("k1 should be the oldest element")
	}
}

func TestLRU_RemaxBytes(t *testing.T) {
	OnEvictedCounter := 0
	OnEvicted := func(_ string, _ Value) {
		OnEvictedCounter++
	}

	cache, _ := NewCache(8, OnEvicted)
	cache.Add(testk1, testv1)
	cache.Add(testk2, testv2)

	// DownmaxBytes
	evicted, _ := cache.RemaxBytes(4)
	if evicted != 1 {
		t.Errorf("1 element should have been evicted: %v", evicted)
	}
	if OnEvictedCounter != 1 {
		t.Errorf("onEvicted should have been called 1 time: %v", OnEvictedCounter)
	}
	if cache.Contains(testk1) {
		t.Errorf("Element k1 should have been evicted")
	}

	cache.Add(testk3, testv3)
	if cache.Contains(testk2) {
		t.Errorf("Element k2 should have been evicted")
	}

	// UpmaxBytes
	evicted, _ = cache.RemaxBytes(12)
	cache.Add(testk1, testv1)
	cache.Add(testk2, testv2)
	if evicted != 0 {
		t.Errorf("0 elements should have been evicted: %v", evicted)
	}
	if k, v, _ := cache.GetOldest(); k != testk3 || v != testv3 {
		t.Errorf("Cache should have k3 elements as oldest element")
	}
}
