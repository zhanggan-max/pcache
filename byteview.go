package pcache

type ByteView struct {
	// 如果 b 不为空，则由 b 存储数据
	b []byte
	// 如果 b 为空，则由 s 存储数据
	s string
}

func (v ByteView) Len() int {
	if v.b != nil {
		return len(v.b)
	}
	return len(v.s)
}

func (v ByteView) ByteSlice() []byte {
	if v.b != nil {
		return cloneBytes(v.b)
	}
	return []byte(v.s)
}

func (v ByteView) String() string {
	if v.b != nil {
		return string(v.b)
	}
	return v.s
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
