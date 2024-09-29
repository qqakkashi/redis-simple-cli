package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
)

func ViewTask(client *redis.Client, taskId string, taskType string) {
	switch taskType {
	case "hash":
		ViewHashTask(client, taskId)
	case "set":
		ViewSetTask(client, taskId)
	case "zset":
		ViewZSetTask(client, taskId)
	case "string":
		ViewStringTask(client, taskId)
	}
}

func ViewHashTask(client *redis.Client, taskID string) {
	currentData, err := client.HGetAll(ctx, taskID).Result()
	if err != nil {
		log.Fatalf("Error getting current hash task data: %s\n", err)
		return
	}

	if len(currentData) == 0 {
		log.Println("Hash task not found.")
		return
	}

	log.Println("Hash Task Data:")
	for field, value := range currentData {
		if isJSON(value) {
			log.Printf("'%s': \n", field)
			prettyPrintJSON(value) // Печатаем отформатированный JSON
		} else {
			log.Printf("'%s': %s\n", field, value)
		}
	}
}

func ViewSetTask(client *redis.Client, taskID string) {
	members, err := client.SMembers(ctx, taskID).Result()
	if err != nil {
		log.Fatalf("Error getting current set task data: %s\n", err)
		return
	}

	if len(members) == 0 {
		log.Println("Set task not found.")
		return
	}

	log.Println("Set Task Data:")
	for _, member := range members {
		log.Printf("Member: %s\n", member)
	}
}

func ViewZSetTask(client *redis.Client, taskID string) {
	members, err := client.ZRangeWithScores(ctx, taskID, 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting current zset task data: %s\n", err)
		return
	}

	if len(members) == 0 {
		log.Println("ZSet task not found.")
		return
	}

	log.Println("ZSet Task Data:")
	for _, member := range members {
		log.Printf("'%s': %f\n", member.Member, member.Score)
	}
}

func ViewStringTask(client *redis.Client, taskID string) {
	value, err := client.Get(ctx, taskID).Result()
	if err != nil {
		log.Fatalf("Error getting current string task data: %s\n", err)
		return
	}

	log.Printf("String Task Data: %s\n", value)
}

func prettyPrintJSON(data string) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, []byte(data), "", "  ")
	if err != nil {
		log.Fatalf("Failed to format JSON: %s\n", err)
		return
	}
	log.Println(prettyJSON.String())
}

func isJSON(str string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(str), &js) == nil
}
