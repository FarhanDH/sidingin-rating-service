package event

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"

	"github.com/FarhanDH/sidingin-rating-service/config"
	"github.com/FarhanDH/sidingin-rating-service/internal/models"
	"github.com/FarhanDH/sidingin-rating-service/internal/service"
	"github.com/rabbitmq/amqp091-go"
)

func failOnOver(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type RabbitMQHandler struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	queue   amqp091.Queue
	svc     *service.RatingService
	logger  *slog.Logger
	cfg     *config.RabbitMQConfig
}

func NewRabbitMQHandler(cfg *config.RabbitMQConfig, svc *service.RatingService, logger *slog.Logger) *RabbitMQHandler {
	return &RabbitMQHandler{
		svc:    svc,
		logger: logger,
		cfg:    cfg,
	}
}

func (c *RabbitMQHandler) Connect(ctx context.Context) error {
	var err error
	c.conn, err = amqp091.Dial(c.cfg.URL)
	if err != nil {
		return err
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		return err
	}

	q, err := c.channel.QueueDeclare(
		c.cfg.Queue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		amqp091.Table{
			"x-queue-type": "quorum", // Key requirement to specify a quorum queue
		}, // arguments
	)
	if err != nil {
		return err
	}
	c.queue = q

	c.logger.Info("RabbitMQ consumer connected",
		slog.String("queue", c.queue.Name),
		slog.String("url", c.cfg.URL),
	)

	return nil
}

func (c *RabbitMQHandler) Consume(ctx context.Context) error {
	msgs, err := c.channel.ConsumeWithContext(
		ctx,
		c.queue.Name,
		"",    // consumer tag
		false, // auto-ack → manual ack for reliability
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	c.logger.Info("RabbitMQ consumer started consuming",
		slog.String("queue", c.queue.Name),
	)

	for msg := range msgs {
		select {
		case <-ctx.Done():
			c.logger.Info("RabbitMQ consumer shutting down",
				slog.String("queue", c.queue.Name),
			)
			return ctx.Err()
		default:
		}

		c.handleMessage(ctx, msg)
	}

	return nil
}

func (c *RabbitMQHandler) handleMessage(ctx context.Context, msg amqp091.Delivery) {
	logger := c.logger.With(
		slog.String("delivery_tag", string(msg.DeliveryTag)),
		slog.String("message_id", msg.MessageId),
		slog.String("correlation_id", msg.CorrelationId),
	)

	var payload models.OrderCompletedPayload
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		logger.Error("failed to unmarshal order completed payload",
			slog.String("error", err.Error()),
			slog.String("body", string(msg.Body)),
		)
		// Don't requeue malformed messages (poison pill)
		if nackErr := msg.Nack(false, false); nackErr != nil {
			logger.Error("failed to nack malformed message", slog.String("error", nackErr.Error()))
		}
		return
	}

	logger.Info("received order completed message",
		slog.String("order_id", payload.OrderID),
		slog.String("customer_id", payload.CustomerID),
		slog.String("technician_id", payload.TechnicianID),
		slog.Time("completed_at", payload.CompletedAt),
	)

	if err := c.svc.HandleCompletedOrder(ctx, payload); err != nil {
		logger.Error("failed to handle completed order",
			slog.String("order_id", payload.OrderID),
			slog.String("error", err.Error()),
		)
		// Requeue for retry on transient failures
		if nackErr := msg.Nack(false, true); nackErr != nil {
			logger.Error("failed to nack message for requeue", slog.String("error", nackErr.Error()))
		}
		return
	}

	logger.Info("order completed message processed successfully",
		slog.String("order_id", payload.OrderID),
	)

	if ackErr := msg.Ack(false); ackErr != nil {
		logger.Error("failed to ack message", slog.String("error", ackErr.Error()))
	}
}

func (c *RabbitMQHandler) Close() error {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			return err
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}
