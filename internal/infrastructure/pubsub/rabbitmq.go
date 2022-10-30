package pubsub

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go.uber.org/zap"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/storyden/internal/config"
)

var _ Bus = &Rabbit{}

type Rabbit struct {
	pub    *amqp.Channel
	sub    *amqp.Channel
	queues map[Topic]amqp.Queue
}

func NewRabbit(cfg config.Config) (Bus, error) {
	conn, err := amqp.Dial(cfg.AmqpAddress)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to connect to amqp server"))
	}

	pub, err := conn.Channel()
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to create publish channel"))
	}

	sub, err := conn.Channel()
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to create subscribe channel"))
	}

	r := Rabbit{pub, sub, make(map[Topic]amqp.Queue)}

	return &r, nil
}

// Declare is meant to be used like:
// const MyQueue string
//
// MyQueue = Declare("name_of_queue")
//
// ps.Publish(MyQueue, ...)
//
func (r *Rabbit) Declare(t string) Topic {
	_, err := r.pub.QueueDeclare(
		t,     // name
		true,  // durable
		false, // auto delete
		false, // exclusive
		false, // nowait
		nil,   // args
	)
	if err != nil {
		panic(err)
	}

	return Topic(t)
}

func (r *Rabbit) Publish(topic Topic, message []byte) error {
	err := r.pub.Publish(
		"",            // exchange
		string(topic), // key
		true,          // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		return fault.Wrap(err, fmsg.With("failed to publish to topic"))
	}

	return nil
}

func (r *Rabbit) Subscribe(topic Topic, handler func([]byte) (bool, error)) error {
	msgs, err := r.sub.Consume(
		string(topic), // queue
		"",            // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		panic(err)
	}

	for m := range msgs {
		ack, err := handler(m.Body)
		if err != nil {
			zap.L().Error("pubsub handler failed", zap.Error(err))
		}

		if ack {
			if err := m.Ack(false); err != nil {
				zap.L().Error("pubsub handler ack failed", zap.Error(err))
			}
		} else {
			if err := m.Nack(false, true); err != nil {
				zap.L().Error("pubsub handler nack failed", zap.Error(err))
			}
		}
	}

	return errors.New("pubsub subscriber hung up")
}
