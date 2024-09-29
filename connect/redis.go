package connect

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"strconv"
)

var ctx = context.Background()

func GetClient() *redis.Client {
	addr := getEnv("REDIS_ADDR", "localhost:6379")
	password := getEnv("REDIS_PASSWORD", "")
	dbNum, err := strconv.Atoi(getEnv("REDIS_DB", "0"))

	if err != nil {
		log.Printf("Incorrect format of REDIS_DB, using 0 by default")
		dbNum = 0
	}

	options := &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbNum,
	}

	client := redis.NewClient(options)

	err = client.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("Can`t connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
	return client
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
