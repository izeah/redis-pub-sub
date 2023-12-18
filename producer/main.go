package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	RegisterProducer(shutdown, f.RedisDB)
}

func RedisInit() {
	redisDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if msg := redisDB.Ping(context.Background()).Val(); msg != "PONG" {
		panic("failed to connect to redis")
	}
	log.Println("successfully connect to redis")
}

func RedisClose() {
	redisDB.Close()
	log.Println("redis connection is closed")
}

func RegisterProducer(shutdown chan os.Signal, redisDB *redis.Client) {
	counter := 0
	fmt.Println("starting producer ...")
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-shutdown:
			fmt.Println("shutdown signal received, exiting ...")
			ticker.Stop()
			return
		case <-ticker.C:
			fmt.Println("rechecking connection ...")
			fmt.Println("result:", redisDB.Ping(context.TODO()).String())
		default:
			if res := redisDB.Publish(context.TODO(), "data_baru", counter); res.Err() != nil {
				fmt.Println("error occured, because:", res.Err().Error())
				continue
			}
			fmt.Println("published to topic:", "data_baru", ", value:", counter)
			counter++
			time.Sleep(3 * time.Second)
		}
	}
}
