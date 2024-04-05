package main

import (
	"bufio"
	"context"
	"fmt"
	"github/justWindy/go_demo/utils"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
	"log"
	"os"
	"strings"
)

func useCyclicBarrier(cli *clientv3.Client, name string) {
	b := recipe.NewBarrier(cli, name)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		action := scanner.Text()
		items := strings.Split(action, " ")
		switch items[0] {
		case "hold":
			b.Hold()
			fmt.Println("hold")
		case "release":
			b.Release()
			fmt.Println("release")
		case "wait":
			b.Wait()
			fmt.Println("after wait")
		case "quit", "exit":
			return
		default:
			fmt.Println("unknown action")
		}

	}
}

func useDoubleCyclicBarrier(cli *clientv3.Client, name string) {
	s1, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()
	b := recipe.NewDoubleBarrier(s1, name, 10)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		action := scanner.Text()
		items := strings.Split(action, " ")
		switch items[0] {
		case "enter":
			b.Enter()
			fmt.Println("enter")
		case "leave":
			b.Leave()
			fmt.Println("leave")
		case "quit", "exit":
			return
		default:
			fmt.Println("unknown action")
		}
	}
}

func doTxnXfer(cli *clientv3.Client, from, to string, amount uint64) (bool, error) {
	getResp, err := cli.Txn(context.Background()).Then(clientv3.OpGet(from), clientv3.OpGet(to)).Commit()
	if err != nil {
		return false, err
	}

	fromKV := getResp.Responses[0].GetResponseRange().Kvs[0]
	toKV := getResp.Responses[1].GetResponseRange().Kvs[0]
	fromV, toV := utils.ToUint64(fromKV.Value), utils.ToUint64(toKV.Value)
	if fromV < amount {
		return false, fmt.Errorf("fromV < amount, insufficient funds")
	}

	txn := cli.Txn(context.Background()).If(
		clientv3.Compare(clientv3.ModRevision(from), "=", fromKV.ModRevision),
		clientv3.Compare(clientv3.ModRevision(to), "=", toKV.ModRevision),
	)

	txn = txn.Then(
		clientv3.OpPut(from, utils.FromUint64(fromV-amount)),
		clientv3.OpPut(to, utils.FromUint64(toV+amount)))

	putResp, err := txn.Commit()
	if err != nil {
		return false, err
	}

	return putResp.Succeeded, nil
}
