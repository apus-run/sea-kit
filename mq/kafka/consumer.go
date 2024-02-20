package kafka

import (
	"context"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"log"
	"os"

	"github.com/IBM/sarama"
	"github.com/apus-run/sea-kit/zlog"
)

type Consumer struct {
	group   sarama.ConsumerGroup
	topics  []string
	groupID string
	handler sarama.ConsumerGroupHandler

	ctx    context.Context
	cancel context.CancelFunc
}

func NewConsumer(logger zlog.Logger, config *sarama.Config, topic string, groupID string, brokers []string, handler *ConsumerGroupHandler) *Consumer {
	sarama.Logger = log.New(os.Stderr, "[sarama_logger]", log.LstdFlags)
	sarama.Logger = logger
	config.Version = sarama.V2_8_2_0

	// Start with a client
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		panic(err)
	}

	// Start a new consumer group
	group, err := sarama.NewConsumerGroupFromClient(groupID, client)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &Consumer{
		group:   group,
		topics:  []string{topic},
		groupID: groupID,
		handler: handler,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (c *Consumer) Consume() error {
	// Track errors
	go func() {
		for err := range c.group.Errors() {
			fmt.Println("ERROR", err)
		}
	}()

	// Iterate over consumer sessions.
	ctx := context.Background()
	for {
		select {
		case <-c.ctx.Done():
			err := c.group.Close()
			logger.Info("[Kafka] Consume ctx done")
			return err
		default:
			if err := c.group.Consume(ctx, c.topics, c.handler); err != nil {
				logger.Errorf("[Kafka] Consume err: %s", err.Error())

				return err
			}
			return nil
		}
	}
}

func (c *Consumer) Close() {
	c.cancel()
}
