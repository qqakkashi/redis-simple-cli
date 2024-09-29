package handlers

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

var ctx = context.Background()

func DeleteTask(client *redis.Client, taskID string, backToList func(client *redis.Client)) {
	err := client.Del(ctx, taskID).Err()
	if err != nil {
		log.Fatalf("Error deleting task: %s", err)
		return
	}
	log.Printf("Task %s deleted successfully\n", taskID)
	backToList(client)
}
