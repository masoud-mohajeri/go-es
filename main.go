package main

import (
	"context"
	"fmt"
	"github.com/masoud-mohajeri/go-es/packages/cron"
	"github.com/masoud-mohajeri/go-es/packages/dispatcher"
	"github.com/masoud-mohajeri/go-es/packages/eventStore"
	"github.com/masoud-mohajeri/go-es/packages/projection"
	"github.com/masoud-mohajeri/go-es/packages/queue"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// INFRASTRUCTURE LAYER
	// keys
	const OrderOne = "ORDER_ONE"

	// topics
	orderOneTopic := queue.NewQueue()

	// dispatchers
	orderOneED := dispatcher.NewDispatcher(orderOneTopic)

	// database
	db := eventStore.NewEventStore()

	// event handlers
	handlers := projection.Handler{
		OrderOne: func(payload string) {
			err := db.Save(OrderOne, []byte(payload))
			if err != nil {
				fmt.Println("-app: error in order one handler")
			}
			state, _ := db.Get(OrderOne)
			fmt.Println("-app: order one current state: \n", string(state), "\n")
		},
	}

	// projection
	p := projection.NewProjection(handlers, []queue.FifoQueue{orderOneTopic})
	go func() {

		err := cron.NewCron(time.Millisecond*500, p.GenerateSnapshot, context.Background())
		if err != nil {
			fmt.Println("-app: error in creating cron job for projection")
		}
	}()

	// BUSINESS LOGIC LAYER ?!
	// order states
	const OrderPlaced = "ORDER_PLACED"
	const OrderPending = "ORDER_PENDING"
	const OrderProcessing = "ORDER_PROCESSING"
	const OrderInTransit = "ORDER_IN_TRANSIT"
	const OrderInDelivery = "ORDER_IN_DELIVERY"
	const OrderDelivered = "ORDER_DELIVERED"
	const OrderReturned = "ORDER_RETURNED"

	// APPLICATION LOGIC LAYER
	orderOneED.Dispatch(OrderPlaced, OrderOne)
	orderOneED.Dispatch(OrderPending, OrderOne)
	orderOneED.Dispatch(OrderProcessing, OrderOne)

	<-time.After(time.Second)
	orderOneED.Dispatch(OrderInTransit, OrderOne)
	<-time.After(time.Second * 2)
	orderOneED.Dispatch(OrderInDelivery, OrderOne)
	<-time.After(time.Second * 3)
	orderOneED.Dispatch(OrderDelivered, OrderOne)
	<-time.After(time.Second * 4)
	orderOneED.Dispatch(OrderReturned, OrderOne)

	// gracefully shutdown
	shutdown := func() {
		//clean-ups
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM)

	select {
	case <-signalCh:
		shutdown()
	}
}
