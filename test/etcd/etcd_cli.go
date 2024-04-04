package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github/justWindy/go_demo/utils"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"os"
	"strings"
)

var (
	nodeID    = flag.Int("id", 0, "node ID")
	addr      = flag.String("addr", "http://192.168.50.89:2379", "etcd address")
	electName = flag.String("name", "my-test-elect", "election name")
)

func main() {
	go func() {
		flag.Parse()
	}()

	// etcd addr
	endpoints := strings.Split(*addr, ",")
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	session, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	e1 := concurrency.NewElection(session, *electName)

	consoleScanner := bufio.NewScanner(os.Stdin)
	for consoleScanner.Scan() {
		action := consoleScanner.Text()
		log.Println("action", action)
		switch action {
		case "elect": // start the elect
			go elect(e1, *electName)
		case "proclaim":
			proclaim(e1, *electName)
		case "resign":
			resign(e1, *electName)
		case "watch":
			go watch(e1, *electName)
		case "query":
			query(e1, *electName)
		case "rev":
			rev(e1, *electName)
		default:
			fmt.Println("unknown action")
		}
	}
}

var count int

func elect(e *concurrency.Election, electName string) {
	log.Println("a campaigning for ID:", *nodeID)
	if err := e.Campaign(context.Background(), fmt.Sprintf("value-%d-%d", *nodeID, count)); err != nil {
		log.Println(err)
	}
	log.Println("campaigned for ID:", *nodeID)
	count++
}

func proclaim(e *concurrency.Election, electName string) {
	log.Println("proclaiming for ID:", *nodeID)
	if err := e.Proclaim(context.Background(), fmt.Sprintf("value-%d-%d", *nodeID, count)); err != nil {
		log.Println(err)
	}
	log.Println("proclaimed for ID:", *nodeID)
	count++
}

func resign(e *concurrency.Election, electName string) {
	log.Println("resigning for ID:", *nodeID)
	if err := e.Resign(context.TODO()); err != nil {
		log.Println(err)
	}
	log.Println("resigned for ID:", *nodeID)
}

func query(e *concurrency.Election, electName string) {
	resp, err := e.Leader(context.Background())
	if err != nil {
		log.Printf("failed to get the current leader: %v", err)
	}
	log.Println("current leader:", utils.Byte2string(resp.Kvs[0].Key), utils.Byte2string(resp.Kvs[0].Value))
}

func rev(e *concurrency.Election, electName string) {
	version := e.Rev()
	log.Println("current rev:", version)
}

func watch(e *concurrency.Election, electName string) {
	ch := e.Observe(context.TODO())

	log.Println("start to watch for ID:", *nodeID)
	for i := 0; i < 10; i++ {
		resp := <-ch
		log.Println("leader changed to", utils.Byte2string(resp.Kvs[0].Key), utils.Byte2string(resp.Kvs[0].Value))
	}
}
