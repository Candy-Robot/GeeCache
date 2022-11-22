// 负责与外部交互，控制缓存存储和获取的主流程
package GeeCache

import (
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// 是组的缓存命名空间，相关的数据会在旁边分布
// 一个 Group 可以认为是一个缓存的命名空间

type Group struct {
	// 每个 Group 拥有一个唯一的名称 name。比如可以创建三个 Group，
	// 缓存学生的成绩命名为 scores， 缓存学生信息的命名为 info， 缓存学生课程的命名为 courses。
	name string
	getter Getter	// 缓存未命中时获取源数据的回调(callback)
	mainCache cache	// 一开始实现的并发缓存
}

var (
	mu		sync.RWMutex
	groups = make(map[string]*Group)
)

// 给组创建一个实例
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil{
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name: name,
		getter: getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// 从键值对的键得到值
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	// 命中的情况
	// 流程 ⑴ 从 mainCache 中查找缓存，如果存在则返回缓存值
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	// 流程 ⑶ ：缓存不存在，则调用 load 方法
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error){
	// 分场景选择，是从本地调用，还是从别的节点调用
	return g.getLocally(key)
}
// 从本地节点寻找key
func (g *Group) getLocally(key string) (ByteView, error) {
	// 通过用户提供的回调函数，获取源数据
	bytes, err := g.getter.Get(key)
	if err != nil { // 如果出了问题之后
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// 将源数据添加到缓存mainCache当中
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

