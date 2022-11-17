package lru

import "container/list"

type Cache struct {
	maxBytes			int64	// 允许使用的最大内存
	nbytes				int64	// 当前已使用的内存
	ll      			*list.List	// 双向链表
	cache			 	map[string]*list.Element	// 键是字符串， 值是双向链表里的指针
	OnEvicted			func(key string, value Value)	//OnEvicted 是某条记录被移除时的回调函数，可以为 nil。
}

// 链表节点的数据类型
type entry struct {
	key string
	value Value
}

// 允许的值是实现了Value接口的任何类型。接口只包含了一个方法。 Len() int，用于返回值所占用的内存大小。
type Value interface {
	Len() int
}

// 方便实例化Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll:list.New(),
		cache: make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 查找功能
// 1、从字典中找到对应的双向链表节点 2、将该节点移动到队尾
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok{
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// 删除
// 实际上是缓存淘汰
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	// 列表不为空
	if ele != nil {
		c.ll.Remove(ele)
		// .(*entry)是进行类型断言 如果是这个类型就不会引发panic 返回T类型的变量
		kv := ele.Value.(*entry)	// 得到key和value
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// 增加或修改
func (c *Cache) Add(key string, value Value) {
	// 修改当前的值
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		// 更新占用的长度
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	}else { // 添加值
		ele := c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	// maxBytes 设置为0表示不对内存大小设限
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}

}