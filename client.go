package pcache

import (
	"go.etcd.io/etcd/clientv3"
)

type client struct {
	name string // 服务名称：pcache/ip:addr
}

func (c *client) Fetch(group string, key string) ([]byte, error) {
	cli, err := clientv3.New(defaultEtcdConfig)
	if err != nil {
		return nil, err
	}
	defer cli.Close()
}
