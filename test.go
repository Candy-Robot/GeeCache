package main

import (
	"fmt"
	"sync"
	"time"
)

var m sync.Mutex
var set = make(map[int]bool, 0)

func printOnce(num int) {
	m.Lock()
	defer m.Unlock()
	if _, exist := set[num]; !exist {
		fmt.Println(num)
	}
	set[num] = true

}

func main() {
	for i := 0; i < 10; i++ {
		// 因为同时读取到这个set中的100了 一开始都是false
		// 应该改为在没写完之前，其他的携程不能读取
		go printOnce(100)
	}
	time.Sleep(time.Second)
}
