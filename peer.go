package pcache

// Picker 定义了节点将请求发送到其他节点的能力
type Picker interface {
	Pick(key string) (Fetcher, bool)
}

// Fetcher 接口定义了向特定客户端请求的能力
type Fetcher interface {
	Fetch(group string, key string) ([]byte, error)
}
