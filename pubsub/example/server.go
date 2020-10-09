package main

import (
	"context"
	"fmt"

	"time"

	"github.com/tx991020/utils/pubsub/broker"
)

func main() {
	fmt.Println("a")

	broker := broker.NewBroker(context.Background(), ":9876", 3)

	go func() {
		m := 0
		for {
			m++
			broker.Publish("iamlegend", []byte(fmt.Sprintf("消息 %d", m)))
			time.Sleep(2 * time.Second)
		}
	}()

	err := broker.Run()
	if err != nil {
		fmt.Println(err)
	}
}
