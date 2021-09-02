package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/kailashyogeshwar85/slim-orderbook/cmd/slim-orderbook/engine"
)

func main() {
	fmt.Println("-----------------------------------------")
	fmt.Println("        SLIM-ORDERBOOK 1.0               ")
	fmt.Println("-----------------------------------------")

	// create consumer for consuming order
	consumer := createConsumer()

	// create producer to send trades and orders data
	producer := createProducer()

	// create the order book
	book := engine.OrderBook{
		Bids: make([]engine.Order, 0, 10000),
		Asks: make([]engine.Order, 0, 10000),
	}

	// create a channel to know when done
	done := make(chan bool)

	// start processing order
	go func() {
		for msg := range consumer.Messages() {
			var order engine.Order
			// deserialize the message
			order.FromJSON(msg.Value)
			// process the order
			trades := book.Process(order)

			log.Println("Trades length: ", len(trades))

			if len(trades) != 0 {
				// send trades to message queue
				for _, trade := range trades {
					rawTrade := trade.ToJSON()

					log.Println("Publishing trade on topic -> trades")
					// publish the message over receiving channel
					producer.Input() <- &sarama.ProducerMessage{
						Topic: "trades",
						Value: sarama.ByteEncoder(rawTrade),
					}
				}
				// mark the offset as commited
				consumer.MarkOffset(msg, "")
			}
		}
		done <- true
	}()
	<-done
}

func createConsumer() *cluster.Consumer {
	// define the configuration for our cluster
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = false
	config.Group.Return.Notifications = false
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // earliest uncommited offset
	config.Consumer.Offsets.CommitInterval = time.Second

	orderTopic := []string{"orders"}

	log.Println("Listening for orders on topic -> ", orderTopic)
	// create the consumer
	consumer, err := cluster.NewConsumer(
		[]string{"127.0.0.1:9092"},
		"orderbook-cg",
		orderTopic,
		config,
	)

	if err != nil {
		log.Fatal("Unable to connect to kafka cluster")
	}

	go handleErrors(consumer)
	go handleNotifications(consumer)
	return consumer
}

func handleErrors(consumer *cluster.Consumer) {
	for err := range consumer.Errors() {
		log.Printf("Error: %s\n", err.Error())
	}
}

func handleNotifications(consumer *cluster.Consumer) {
	for ntf := range consumer.Notifications() {
		log.Printf("Rebalanced %+v\n", ntf)
	}
}

func createProducer() sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = false         // fire and forget
	config.Producer.Return.Errors = true             // notify on failed
	config.Producer.RequiredAcks = sarama.WaitForAll // waits for all insync replicas to commit

	producer, err := sarama.NewAsyncProducer([]string{"127.0.0.1:9092"}, config)

	if err != nil {
		log.Fatal("Unable to connect producer to kafka server")
	}

	return producer

}
