package main

import (
	"fmt"
	"log"
	"pcache"
	"sync"
)

var mysql = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func GetFromMysql(key string) ([]byte, error) {
	value, ok := mysql[key]
	if !ok {
		return []byte{}, fmt.Errorf("no key %v", key)
	}
	return []byte(value), nil
}

func GetScore(group *pcache.Group, key string, wg *sync.WaitGroup) {
	defer wg.Done()
	view, err := group.Get(key)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(view.String())
}

func main() {
	group := pcache.NewGroup("score", 2<<10, pcache.GetterFunc(GetFromMysql))
	addr := "localhost:9999"
	server, err := pcache.NewServer(addr)
	if err != nil {
		log.Fatal(err)
	}
	server.SetPeers(addr)
	group.RegisterPicker(server)
	log.Println("pcache is running at: ", addr)
	go func() {
		err = server.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(4)
	go GetScore(group, "Tom", &wg)
	go GetScore(group, "Tom", &wg)
	go GetScore(group, "Tom", &wg)
	go GetScore(group, "Tom", &wg)
	wg.Wait()
}
