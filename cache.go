// 并发控制
package GeeCache

import (
	"lru"
	"sync"
)

type cache struct {
	mu 		sync.Mutex
	lru		*lru.Cache
	cacheBytes	int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		// lru 缓存淘汰机制
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil{
		return
	}
	if v, ok := c.lru.Get(key); ok{
		// (ByteView)是进行类型断言 如果是这个类型就不会引发panic 返回T类型的变量
		return v.(ByteView), ok
	}

	return
}