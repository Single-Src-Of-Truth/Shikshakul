package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type EventPublisher interface {
	PublishEmailEvent(ctx context.Context, event EmailEvent) error
}

type redisPublisher struct {
	client *redis.Client
	logger *zap.Logger
}

func NewRedisPublisher(client *redis.Client, logger *zap.Logger) EventPublisher {
	return &redisPublisher{
		client: client,
		logger: logger,
	}
}

func (p *redisPublisher) PublishEmailEvent(ctx context.Context, event EmailEvent) error {
	streamName := StreamEmailLow
	if event.Priority == PriorityHigh {
		streamName = StreamEmailHigh
	}

	payloadBytes, err := json.Marshal(event)
	if err != nil {
		p.logger.Error("Failed to marshal email event", zap.Error(err), zap.String("event_id", event.EventID))
		return fmt.Errorf("failed to marshal email event: %w", err)
	}

	err = p.client.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: map[string]interface{}{
			"payload": string(payloadBytes),
		},
	}).Err()

	if err != nil {
		p.logger.Error("Failed to publish email event to Redis",
			zap.Error(err),
			zap.String("stream", streamName),
			zap.String("event_id", event.EventID),
			zap.String("source", event.Source),
		)
		return err
	}

	p.logger.Info("Email event published successfully",
		zap.String("stream", streamName),
		zap.String("event_id", event.EventID),
		zap.String("source", event.Source),
	)
	return nil
}
