// 缓存值的抽象与封装
package GeeCache

// 是一个不可见变的字节视图
type ByteView struct {
	b []byte
}

// 返回的是视图的长度
func (v ByteView) Len() int {
	return len(v.b)
}

// 将数据的副本作为字节切片返回 防止被外部程序修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// 将数据作为字符串返回
func (v ByteView) String() string{
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}


