package main

import (
	"context"
	"flag"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"math/rand"
	"strings"
	"sync"
)

var (
	addr = flag.String("addr", "http://192.168.50.89:12379", "etcd address")
)

func main() {
	flag.Parse()

	endpoints := strings.Split(*addr, ",")

	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	totalAccounts := 5
	for i := 0; i < totalAccounts; i++ {
		k := fmt.Sprintf("account-%d", i)
		if _, err = cli.Put(context.Background(), k, "100"); err != nil {
			log.Fatal(err)
		}
	}

	// STM apply func
	exchange := func(stm concurrency.STM) error {
		from, to := rand.Intn(totalAccounts), rand.Intn(totalAccounts)
		if from == to {
			return nil
		}

		fromK, toK := fmt.Sprintf("account-%d", from), fmt.Sprintf("account-%d", to)
		fromV, toV := stm.Get(fromK), stm.Get(toK)
		fromInt, toInt := 0, 0
		fmt.Sscanf(fromV, "%d", &fromInt)
		fmt.Sscanf(toV, "%d", &toInt)

		xfer := fromInt / 2
		fromInt, toInt = fromInt-xfer, toInt+xfer

		stm.Put(fromK, fmt.Sprintf("%d", fromInt))
		stm.Put(toK, fmt.Sprintf("%d", toInt))
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				if _, serr := concurrency.NewSTM(cli, exchange); serr != nil {
					log.Fatal(err)
				}
			}
		}()
	}

	wg.Wait()

	sum := 0
	accts, err := cli.Get(context.Background(), "account-", clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}
	for _, kv := range accts.Kvs {
		v := 0
		fmt.Sscanf(string(kv.Value), "%d", &v)
		sum += v
		log.Printf("account %s: %d", kv.Key, v)
	}

	log.Println("account sum is ", sum)
}
