package utils

import "time"

var lru *LruCache

func init() {
	lru = NewLruCache(10000, func(key, value interface{}) {})
	lru.AddTimeoutAfterCreate(10 * time.Minute)
	lru.AddTimeoutAfterRead(10 * time.Minute)
}

func CachePut(key, value interface{}) {
	lru.Add(key, value)
}

func CacheGet(key interface{}) (interface{}, bool) {
	return lru.Get(key)
}

func CacheRemove(key interface{}) bool {
	return lru.Remove(key)
}
