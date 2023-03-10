package rabbitmq

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

//go:generate mockery --name AMQPPublisher --case underscore --output mocks/

// AMQPPublisher ...
type AMQPPublisher interface {
	PublishWithContext(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	Publish(ctx context.Context, exchange, key string, msg amqp.Publishing) error
	GetChannel() *amqp.Channel
	GetConnection() *amqp.Connection
	Close()
}

// Publisher ...
type Publisher struct {
	AMQPConn *amqp.Connection
	AMQPChan *amqp.Channel
}

// NewPublisher ...
func NewPublisher(cfg *Config) (*Publisher, error) {
	conn, err := NewRabbitMQConnection(cfg)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Publisher{AMQPConn: conn, AMQPChan: channel}, nil
}

// NewRabbitMQConnection ...
func NewRabbitMQConnection(cfg *Config) (*amqp.Connection, error) {
	conn, err := amqp.Dial(cfg.URI)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Close ...
func (p *Publisher) Close() {
	log.Println("Closing amqp channel.")
	if err := p.AMQPChan.Close(); err != nil {
		log.Printf("publisher mqpChan.Close err: %v", err)
	}

	log.Println("Closing amqp connection.")
	if err := p.AMQPConn.Close(); err != nil {
		log.Printf("publisher amqpConn.Close err: %v", err)
	}
}

// Publish ...
func (p *Publisher) Publish(ctx context.Context, exchange, key string, msg amqp.Publishing) error {
	if err := p.AMQPChan.PublishWithContext(ctx, exchange, key, false, false, msg); err != nil {
		log.Printf("publisher Publish err: %v", err)
		return err
	}
	return nil
}

// PublishWithContext ...
func (p *Publisher) PublishWithContext(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	if err := p.AMQPChan.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg); err != nil {
		log.Printf("publisher PublishWithContext err: %v", err)
		return err
	}
	return nil
}

func (p *Publisher) GetChannel() *amqp.Channel {
	return p.AMQPChan
}

func (p *Publisher) GetConnection() *amqp.Connection {
	return p.AMQPConn
}
