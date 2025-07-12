package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisHelper struct {
	client      *redis.Client
	channel     string
	id          int
	heartbeatT  *time.Timer
	heartbeatCh chan struct{}
}

func NewRedisHelper(channel, addr string, id int) *RedisHelper {
	rdb := redis.NewClient(&redis.Options{Addr: addr})

	r := &RedisHelper{
		client:      rdb,
		channel:     channel,
		id:          id,
		heartbeatT:  time.NewTimer(6 * time.Second),
		heartbeatCh: make(chan struct{}, 1),
	}
	go r.monitorTimeout()
	return r
}

func (r *RedisHelper) Publish(channel, msg string) {
	r.client.Publish(context.Background(), channel, msg)
}

func (r *RedisHelper) Subscribe(channel string, handle func(string)) {
	sub := r.client.Subscribe(context.Background(), channel)
	ch := sub.Channel()
	go func() {
		for msg := range ch {
			handle(msg.Payload)
		}
	}()
}

func (r *RedisHelper) monitorTimeout() {
	for {
		<-r.heartbeatT.C
		select {
		case r.heartbeatCh <- struct{}{}:
		default:
		}
		r.heartbeatT.Reset(6 * time.Second)
	}
}

func (r *RedisHelper) RefreshHeartbeatTimer() {
	if !r.heartbeatT.Stop() {
		select {
		case <-r.heartbeatT.C:
		default:
		}
	}
	r.heartbeatT.Reset(6 * time.Second)
}

func (r *RedisHelper) HeartbeatTimeout() <-chan struct{} {
	return r.heartbeatCh
}
