package go_mq

import (
	"context"
	"fmt"
	"time"
)

// Delivery Functional interfaces for message distribution
type Delivery interface {
	// Payload Content of the transmitted message
	Payload() string

	Ack() error
	Reject() error
	Push() error
}

type redisDelivery struct {
	ctx context.Context

	payload     string
	unackedKey  string
	rejectedKey string
	pushKey     string

	redisClient RedisClient

	errChan chan<- error
}

func NewDelivery(
	ctx context.Context,
	payload string,
	unackedKey string,
	rejectedKey string,
	pushKey string,
	redisClient RedisClient,
	errChan chan<- error,
) *redisDelivery {
	return &redisDelivery{
		ctx:         ctx,
		payload:     payload,
		unackedKey:  unackedKey,
		rejectedKey: rejectedKey,
		pushKey:     pushKey,
		redisClient: redisClient,
		errChan:     errChan,
	}
}

func (d *redisDelivery) String() string {
	return fmt.Sprintf("[%s %s]", d.payload, d.unackedKey)
}

func (d *redisDelivery) Payload() string {
	return d.payload
}

func (d *redisDelivery) Ack() error {
	var (
		errCount = 0
	)

	for {
		count, err := d.redisClient.LRem(d.unackedKey, 1, d.payload)
		// if no redis error
		if err == nil {
			if count == 0 {
				return ErrorNotFound
			}
			return nil
		}

		// redis error
		errCount++
		select {
		case d.errChan <- &DeliveryError{Delivery: d, RedisErr: err, Count: errCount}:
		default:
		}

		if err := d.ctx.Err(); err != nil {
			return ErrorConsumingStopped
		}

		// a bad way, just for test case
		time.Sleep(time.Second)
	}
}

func (d *redisDelivery) Reject() error {
	return d.move(d.rejectedKey)
}

func (d *redisDelivery) Push() error {
	if d.pushKey == "" {
		return d.Reject() // fall back to rejecting
	}

	return d.move(d.pushKey)
}

func (d *redisDelivery) move(key string) error {
	var (
		errCount = 0
	)
	for {
		_, err := d.redisClient.LPush(key, d.payload)
		if err == nil {
			break
		}
		errCount++
		select {
		case d.errChan <- &DeliveryError{Delivery: d, RedisErr: err, Count: errCount}:
		default:
		}

		if err := d.ctx.Err(); err != nil {
			return ErrorConsumingStopped
		}

		time.Sleep(time.Second)
	}

	return d.Ack()
}
