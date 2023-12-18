package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"
)

var redisDB *redis.Client

type Factory struct {
	RedisDB *redis.Client
}

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	signal.Notify(shutdown, syscall.SIGTERM)

	RedisInit()
	defer RedisClose()

	f := &Factory{
		RedisDB: redisDB,
	}

	RegisterConsumer(shutdown, f.RedisDB)
}

func RedisInit() {
	redisDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if msg := redisDB.Ping(context.TODO()).Val(); msg != "PONG" {
		panic("failed to connect to redis")
	}
	log.Println("successfully connect to redis")
}

func RedisClose() {
	redisDB.Close()
	log.Println("redis connection is closed")
}

func RegisterConsumer(shutdown chan os.Signal, redisDB *redis.Client) {
	fmt.Println("starting Consumer ...")
	pubsub := redisDB.Subscribe(context.TODO(), "data_baru")
	if _, err := pubsub.Receive(context.TODO()); err != nil {
		panic("failed to connect to redis subscriber")
	}
	ch := pubsub.Channel()
	for {
		select {
		case msg := <-ch:
			if msg != nil {
				fmt.Println("consumed data:", msg.Payload)
			}
		case <-shutdown:
			_ = pubsub.Close()
			fmt.Println("Pubsub closed. Channel closed.")
			fmt.Println("shutdown signal received, exiting ...")
			return
		}
	}
}
