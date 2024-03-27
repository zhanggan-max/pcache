package pcache

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"pcache/consistenthash"
	pb "pcache/pcachepb"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

const (
	defaultAddr    = "127.0.0.1:6324"
	defaultRepicas = 50
)

var defaultEtcdConfig = clientv3.Config{
	Endpoints:   []string{"localhost:2379"},
	DialTimeout: 5 * time.Second,
}

type server struct {
	pb.UnimplementedPcacheServer

	addr           string // ip:port
	status         bool   //true: running false: stop
	stopSignal     chan error
	mu             sync.Mutex
	consistentHash *consistenthash.Map
	clients        map[string]*client
}

func NewServer(addr string) (*server, error) {
	if addr == "" {
		addr = defaultAddr
	}
	return &server{addr: addr}, nil
}

// Get 是 rpc 服务要求的方法
func (s *server) Get(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	group, key := in.GetGroup(), in.GetKey()
	repv := &pb.Response{}

	log.Printf("[pcache server %s] Recv RPC Request - (%s)/(%s)", s.addr, group, key)
	if key == "" {
		return repv, fmt.Errorf("key required")
	}
	g := GetGroup(group)
	if g == nil {
		return repv, fmt.Errorf("group not found")
	}
	view, err := g.Get(key)
	if err != nil {
		return repv, err
	}
	repv.Value = view.ByteSlice()
	return repv, nil
}

// Start 启动服务器，todo: 并将其注册到 etcd 中
func (s *server) Start() error {
	s.mu.Lock()
	if s.status {
		s.mu.Unlock()
		return fmt.Errorf("server already started")
	}
	s.status = true
	s.stopSignal = make(chan error)

	port := strings.Split(s.addr, ":")[1]
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterPcacheServer(grpcServer, s)

	//todo: etcd registe
	s.mu.Unlock()
	if err := grpcServer.Serve(lis); s.status && err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}

// SetPeers 方法将服务实例注册到 Server 中
func (s *server) SetPeers(peerAddrs ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.consistentHash = consistenthash.New(defaultRepicas, nil)
	s.consistentHash.Registe(peerAddrs...)
	s.clients = make(map[string]*client)
	for _, peerAddr := range peerAddrs {
		service := fmt.Sprintf("pcache%s", peerAddr)
		// todo: 实现 newclient 方法
		s.clients[peerAddr] = NewClient(service)
	}
}

// Pick 使用一致性哈希算法选择 key 应使用的 cache
// false 表示从本地获取
func (s *server) Pick(key string) (Fetcher, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	peerAddr := s.consistentHash.GetPeer(key)
	if peerAddr == s.addr {
		log.Printf("pick self %s\n", s.addr)
		return nil, false
	}
	log.Printf("cache %s pick remote peer: %s\n", s.addr, peerAddr)
	return s.clients[peerAddr], true
}

// 停止 server 运行
func (s *server) Stop() {
	s.mu.Lock()
	if !s.status {
		s.mu.Unlock()
		return
	}
	s.stopSignal <- nil // 停止发送 KeepAlive 信号
	s.status = false    // 设置服务状态为 stop
	s.clients = nil
	s.consistentHash = nil
	s.mu.Unlock()
}

// 要求 Server 实现 Picker 接口
var _ Picker = (*server)(nil)
