package main

import (
	"context"
	"fmt"
	"time"

	gofakeIt "github.com/brianvoe/gofakeit/v6"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/rodkevich/mvpbe/pkg/rabbitmq"
)

// docker run --hostname my-rabbit --name some-rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3.9-management
// http://localhost:15672 for ui
func main() {
	cfg := rabbitmq.Config{URI: "amqp://guest:guest@localhost:5672"}
	conn, err := NewRabbitMQConnection(cfg)
	if err != nil {
		println("NewRabbitMQConnection ", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		println("conn.Channel ", err)
	}

	err = channel.ExchangeDeclare("sample_exchange", "direct", true, false, false, false, nil)
	if err != nil {
		println("channel.ExchangeDeclare ", err)
	}

	q, err := channel.QueueDeclare("sample_que", true, false, false, false, nil)
	if err != nil {
		println("channel.QueueDeclare ", err)
	}
	fmt.Println(q)
	p := publisher{
		amqpConn: conn,
		amqpChan: channel,
	}
	defer p.Close()

	ctx := context.Background()
	err = p.Publish(ctx, "", "sample_que", amqp.Publishing{
		Headers:   map[string]interface{}{"trace-id": gofakeIt.UUID()},
		Timestamp: time.Now().UTC(),
		Body:      []byte(gofakeIt.BeerAlcohol()),
	},
	)

	if err != nil {
		println("p.Publish ", err)
	}
	err = p.Publish(ctx, "sample_exchange", "",
		amqp.Publishing{
			Headers:   map[string]interface{}{"trace-id": gofakeIt.UUID()},
			Timestamp: time.Now().UTC(),
			Body:      []byte(gofakeIt.BeerAlcohol()),
		},
	)

	if err != nil {
		println("p.Publish ", err)
	}
}

func NewRabbitMQConnection(cfg rabbitmq.Config) (*amqp.Connection, error) {
	conn, err := amqp.Dial(cfg.URI)
	if err != nil {
		println(err)
		return nil, err
	}
	return conn, nil
}

type publisher struct {
	amqpConn *amqp.Connection
	amqpChan *amqp.Channel
}

func (p *publisher) Publish(ctx context.Context, exchange, key string, msg amqp.Publishing) error {
	if err := p.amqpChan.PublishWithContext(ctx, exchange, key, false, false, msg); err != nil {
		err = fmt.Errorf("publisher Publish err: %w", err)
		println("amqpChan.PublishWithContext ", err.Error())

		return err
	}
	return nil
}

func (p *publisher) Close() error {
	if err := p.amqpChan.Close(); err != nil {
		e := fmt.Errorf("publisher mqpChan.Close err: %w", err)
		println("amqpChan.Close ", err.Error())

		return e
	}
	if err := p.amqpConn.Close(); err != nil {
		e := fmt.Errorf("publisher amqpConn.Close err: %w", err)
		println("amqpConn.Close ", err.Error())

		return e
	}
	return nil
}
