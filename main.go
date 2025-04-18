package main

import (
	"os"
	"log"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	fmt.Println(os.Getenv("REDIS_ADDR"))
	fmt.Println(os.Getenv("REDIS_PASSWORD"))

	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
			log.Fatal("Redis connect error:", err)
	}
	fmt.Println("Redis:", pong)

	InitTelegramAPI()
}
