// Copyright (c) 2019. stc Inc. All rights reserved.

package main

import (
	"context"
	"flag"
	"fmt"

	"log"
	"os"
	"os/signal"

	"github.com/tx991020/utils/pubsub"
	"github.com/tx991020/utils/pubsub/client"
)

type DemoQueryFilter struct {
	Address string
}

func (d *DemoQueryFilter) GetTopic() string {
	return d.Address
}

var addr = flag.String("addr", "localhost:9876", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	client, err := client.Dial(*addr)
	if err != nil {
		log.Fatal(err)
	}

	query := &DemoQueryFilter{"iamlegend"}
	logs := make(chan pubsub.SubMessage)

	sub, err := client.Subscribe(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(sub)

	for {
		select {
		case <-interrupt:
			return

		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Println("你好，正式开始读数据:")
			fmt.Println(string(vLog.GetConent())) // pointer to event log
		}
	}

	fmt.Println(client)

}
