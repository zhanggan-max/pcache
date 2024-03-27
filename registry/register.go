package registry

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/clientv3"
)

var defaultEtcdConfig = clientv3.Config{
	Endpoints:   []string{"localhost:2397"},
	DialTiemout: 5 * time.Second,
}

// etcdAdd 在租赁模式添加一对 kv 至 etcd
func etcdAdd(c *clientv3.Client, lid clientv3.LeaseID, service string, addr string) error {
	em, err := endpoints.NewManager(c, service)
	if err != nil {
		return err
	}
	return em.AddEndPoint(c.Ctx(), service+"/"+addr, endpoints.Endpoint{Addr: addr}, clientv3.WithLease(lid))
}

// Register 注册一个服务到 etcd
// 如果不出错，Register 不会返回
func Register(service string, addr string, stop chan error) error {
	cli, err := clientv3.New(defaultEtcdConfig)
	if err != nil {
		return fmt.Errorf("create etcd client failed: %v", err)
	}
	defer cli.Close()
	resp, err := cli.Grant(context.Background(), 5)
	if err != nil {
		return fmt.Errorf("create lease failed: %v", err)
	}
	leaseID := resp.ID
	err = etcdAdd(cli, leaseID, service, addr)
	if err != nil {
		return fmt.Errorf("add etcd record failed: %v", err)
	}
	ch, err := cli.KeepAlive(context.Background(), leaseID)
	if err != nil {
		return fmt.Errorf("set keepalive failed: %v", err)
	}
	log.Printf("%s register service ok\n", addr)
	for {
		select {
		case err := <-stop:
			if err != nil {
				log.Println(err)
			}
			return err
		case <-cli.Ctx().Done():
			log.Println("service closed")
			return nil
		case _, ok := <-ch:
			if !ok {
				log.Println("keep alive channel closed")
				_, err := cli.Revoke(context.Background(), leaseID)
				return err
			}
		}
	}
}
