package main

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	getenv := os.Getenv("AMQP_URI")

	conn, err := amqp.Dial(getenv)
	if err != nil {
		log.Fatal(err)
	}
	println("connected to:", getenv)

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [que_name]...", os.Args[0])
		_ = conn.Close()
		os.Exit(1)
	}

	msgs, err := ch.Consume(
		os.Args[1], // que name
		"cli",      // consumer
		true,       // auto ack
		false,      // exclusive
		false,      // non local
		false,      // no wait
		nil,        // args
	)
	if err != nil {
		_ = conn.Close()
		log.Fatal(err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("[+] %s", d.Body)
		}
	}()

	<-forever
}
