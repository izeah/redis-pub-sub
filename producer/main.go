package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	signal.Notify(shutdown, syscall.SIGTERM)

	RegisterProducer(shutdown, RedisClient)

	<-shutdown

	fmt.Println("shutdown")
}

func RegisterProducer(shutdown chan os.Signal, redisDB *redis.Client) {
	counter := 0
	fmt.Println("starting producer ...")
	ticker := time.NewTicker(60 * time.Second)
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
