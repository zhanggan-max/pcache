package pcache

import (
	"context"
	"fmt"
	pb "pcache/pcachepb"
	"time"

	"google.golang.org/grpc"
)

type client struct {
	name string // 服务名称：pcache/ip:addr
}

// Fetch 从 remote peer 获取对应的缓存值
func (c *client) Fetch(group string, key string) ([]byte, error) {
	conn, err := grpc.Dial(c.name)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	grpcClient := pb.NewPcacheClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := grpcClient.Get(ctx, &pb.Request{Group: group, Key: key})
	if err != nil {
		return nil, fmt.Errorf("could not get %s/%s from peer %s", group, key, c.name)
	}
	return resp.GetValue(), nil
}

func NewClient(service string) *client {
	return &client{name: service}
}

var _ Fetcher = (*client)(nil)
