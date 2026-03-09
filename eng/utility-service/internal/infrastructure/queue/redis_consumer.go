package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/utility-service/internal/service"
	"github.com/Single-Src-Of-Truth/Shikshakul/lib/core-go/events"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type StreamConsumer struct {
	client       *redis.Client
	emailService *service.EmailService
	logger       *zap.Logger
}

func NewStreamConsumer(client *redis.Client, emailService *service.EmailService, logger *zap.Logger) *StreamConsumer {
	return &StreamConsumer{
		client:       client,
		emailService: emailService,
		logger:       logger,
	}
}

func (c *StreamConsumer) Start(ctx context.Context, streamName string, groupName string, workerName string) {
	c.logger.Info("Starting Redis Stream Consumer", zap.String("stream", streamName), zap.String("group", groupName))

	err := c.client.XGroupCreateMkStream(ctx, streamName, groupName, "$").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		c.logger.Fatal("Failed to create consumer group", zap.Error(err))
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				c.logger.Info("Stopping consumer loop", zap.String("worker", workerName))
				return
			default:
				streams, err := c.client.XReadGroup(ctx, &redis.XReadGroupArgs{
					Group:    groupName,
					Consumer: workerName,
					Streams:  []string{streamName, ">"},
					Count:    10,
					Block:    2 * time.Second,
				}).Result()

				if err != nil {
					if err == redis.Nil {
						continue
					}
					c.logger.Error("Error reading from stream", zap.Error(err))
					time.Sleep(1 * time.Second)
					continue
				}

				for _, stream := range streams {
					for _, msg := range stream.Messages {
						c.processMessage(ctx, streamName, groupName, msg)
					}
				}
			}
		}
	}()
}

func (c *StreamConsumer) processMessage(ctx context.Context, streamName string, groupName string, msg redis.XMessage) {
	payloadStr, ok := msg.Values["payload"].(string)
	if !ok {
		c.logger.Error("Malformed event in stream: missing 'payload' key", zap.String("msg_id", msg.ID))
		return
	}

	var event events.EmailEvent
	if err := json.Unmarshal([]byte(payloadStr), &event); err != nil {
		c.logger.Error("Failed to parse event JSON", zap.Error(err), zap.String("payload", payloadStr))
		c.client.XAck(ctx, streamName, groupName, msg.ID)
		return
	}

	err := c.emailService.ProcessEmailEvent(ctx, event)
	if err != nil {
		c.logger.Error("Failed to process email event", zap.String("event_id", event.EventID), zap.Error(err))
		return
	}

	c.client.XAck(ctx, streamName, groupName, msg.ID)
}
