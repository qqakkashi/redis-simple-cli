package helpers

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

var ctx = context.Background()

func GetTypeByTaskId(client *redis.Client, taskId string) string {
	taskType, err := client.Type(ctx, taskId).Result()
	if err != nil {
		log.Fatalf("Error while getting type for task-%s. Error:%s \n", taskId, err)
	}
	return taskType
}
