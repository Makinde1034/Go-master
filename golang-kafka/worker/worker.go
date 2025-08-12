package main

import (
	
	"github.com/IBM/sarama"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	topic := "comments"
	worker, err := connectConsumer([]string{"localhost:9092"})
	if err != nil {
		println("Failed to connect to Kafka producer:", err.Error())
		panic(err)
	}

	consumer, err := worker.ConsumePartition(topic, 0,sarama.OffsetNewest)

	if err != nil {
		println("Failed to connect to Kafka producer:", err.Error())
		panic(err)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	msgCount := 0

	doneChan := make(chan struct{})

	go func(){
		for {
			select {
				case msg := <-consumer.Messages():
					println("Received message:", string(msg.Value))
					msgCount++
				case <-sigchan:
					println("Interrupt signal received, shutting down...")
					doneChan <- struct{}{}
					return
				}
		}
	}()

	<- doneChan
	if err := worker.Close(); err != nil {
		panic(err)
		println("Failed to close consumer:", err.Error())
	}
}


func connectConsumer(brokerUrl []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = false
	conn, err := sarama.NewConsumer(brokerUrl, config)	
	if err != nil {
		println("Error connecting to Kafka:", err.Error())
		return nil, err
	}

	return conn, nil

}