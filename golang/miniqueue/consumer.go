package miniqueue

import (
	"context"
	"errors"
	"fmt"
)

type eventType int

type serverError string

var (
	rrRequestCancelled  = serverError("request context cancelled")
)
const (
	eventTypePublish eventType = iota
	eventTypeNack
	eventTypeBack
)

type notifier interface {
	NotifyConsumer(topic string, ev eventType)
}

// consumers handles providing values iteratively to a single consumers.
type consumer struct {
	id string
	topic string
	ackOffset int
	store storer
	eventChan chan eventType
	notifier notifier
}

func (c *consumer) Next(ctx context.Context) (val value, err error) {
	val, ao, err := c.store.GetNext(c.topic)
	if errors.Is(err, errTopicEmpty) {
		select {
		case <-c.eventChan:
		case <-ctx.Done():
			return nil, errRequestCancelled
		}
		return c.Next(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("getting next from store: %v", err)
	}

	c.ackOffset = ao

	return val, err
}

func (c *consumer) Ack() error {
	if err := c.store.Ack(c.topic, c.ackOffset); err != nil {
		return fmt.Errorf("acking topic %s with offset %d: %v", c.topic, c.ackOffset, err)
	}

	return nil
}

func (c *consumer) Nack() error {
	if err := c.store.Nack(c.topic, c.ackOffset); err != nil {
		return fmt.Errorf("nacking topic %s with offset %d: %v", c.topic, c.ackOffset, err)
	}

	c.notifier.NotifyConsumer(c.topic, eventTypeNack)

	return nil
}

func (c *consumer) Back() error {
	if err := c.store.Back(c.topic, c.ackOffset); err != nil {
		return fmt.Errorf("nacking topic %s with offset %d: %v", c.topic, c.ackOffset, err)
	}

	c.notifier.NotifyConsumer(c.topic, eventTypeBack)

	return nil
}

// EventChan returns a channel to notify the consumers of events occurring on the topic
func (c *consumer) EventChan() <-chan eventType {
	return c.eventChan
}