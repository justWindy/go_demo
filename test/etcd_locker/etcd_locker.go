package main

import (
	"context"
	"flag"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"math/rand"
	"strings"
	"time"
)

var (
	addr     = flag.String("addr", "http://192.168.50.89:12379", "etcd address")
	lockName = flag.String("name", "lock", "lock name")
)

func main() {
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	endpoints := strings.Split(*addr, ",")

	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	useLock(cli)

	useMutex(cli)
}

func useLock(cli *clientv3.Client) {
	// generate the session for the locker
	s1, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()
	locker := concurrency.NewLocker(s1, *lockName)

	// required locker
	log.Println("acquiring lock")
	locker.Lock()
	log.Println("acquired lock")

	//wait for a moment
	time.Sleep(time.Duration(rand.Intn(30)) * time.Second)
	locker.Unlock()

	// release locker
	log.Println("released lock")
}

func useMutex(cli *clientv3.Client) {
	// generate the session for the mutex
	s1, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()
	locker := concurrency.NewMutex(s1, *lockName)

	// before acquire the lock, query it
	log.Printf("before acquiring lock, key: %s", locker.Key())
	//acquire lock
	log.Println("acquiring lock")
	if err := locker.Lock(context.Background()); err != nil {
		log.Fatal(err)
	}
	log.Printf("acquired lock, key: %s", locker.Key())

	// wait for a moment
	time.Sleep(time.Duration(rand.Intn(30)) * time.Second)

	//release lock
	if err := locker.Unlock(context.Background()); err != nil {
		log.Fatal(err)
	}

	log.Println("released lock")
}
