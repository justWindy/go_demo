package main

import (
	"bufio"
	"flag"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	addr      = flag.String("addr", "127.0.0.1:2379", "etcd address")
	queueName = flag.String("queueName", "my-test-queue", "etcd queue name")
)

func main() {
	flag.Parse()

	endpoints := strings.Split(*addr, ",")
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	q := recipe.NewQueue(cli, *queueName)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		action := scanner.Text()
		items := strings.Split(action, " ")
		switch items[0] {
		case "push":
			if len(items) != 2 {
				fmt.Println("invalid add action, must set value to push")
				continue
			}
			q.Enqueue(items[1])
		case "pop":
			v, err := q.Dequeue()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("pop:", v)
		case "quit", "exit":
			return
		default:
			fmt.Println("unknown action:", items[0])
		}
	}
}

func usePriorityQueue(cli *clientv3.Client, name string) {
	q := recipe.NewPriorityQueue(cli, name)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		action := scanner.Text()
		items := strings.Split(action, " ")
		switch items[0] {
		case "push":
			if len(items) != 3 {
				fmt.Println("invalid add action, must set value to push")
				continue
			}
			pr, err := strconv.Atoi(items[2])
			if err != nil {
				log.Println("must set uint16 as priority")
				continue
			}
			q.Enqueue(items[1], uint16(pr))
		case "pop":
			v, err := q.Dequeue()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("pop:", v)
		case "quit", "exit":
			return
		default:
			fmt.Println("unknown action:", items[0])
		}
	}
}
