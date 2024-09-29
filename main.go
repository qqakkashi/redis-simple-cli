package main

import (
	cliui "github.com/qqakkashi/redis-simple-cli/cli-ui"
	"github.com/qqakkashi/redis-simple-cli/connect"
	"github.com/redis/go-redis/v9"
	"log"
)

func main() {
	client := connect.GetClient()
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			log.Println("Can`t close the connection error:", err)
		}
	}(client)

	for {
		action := cliui.DisplayMenu()
		switch action {
		case 0:
			cliui.DisplayTasks(client)
		case 1:
			return
		default:
			log.Println("Invalid choice. Please try again.")
		}
	}
}
