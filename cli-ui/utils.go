package cli_ui

import (
	"context"
	"errors"
	"github.com/manifoldco/promptui"
	"github.com/qqakkashi/redis-simple-cli/handlers"
	"github.com/qqakkashi/redis-simple-cli/helpers"
	"github.com/redis/go-redis/v9"
	"log"
)

var ctx = context.Background()

func DisplayMenu() int {
	actions := []string{"Go to tasks", "Exit"}
	action, _ := GetPrompt("Main menu", actions)
	return action
}

func DisplayTasks(redis *redis.Client) {
	tasks, err := redis.Keys(ctx, "*").Result()
	if err != nil {
		log.Fatalf("Error while getting tasks: %s\n", err)
	}
	action, taskId := GetPrompt("Select a task", tasks)
	if action == -1 {
		return
	}
	taskType := helpers.GetTypeByTaskId(redis, taskId)
	ManageTaskMenu(redis, taskId, taskType)
}

func ManageTaskMenu(redis *redis.Client, taskId string, taskType string) int {
	for {
		actions := []string{"View task info", "Edit task", "Delete task", "Back to list"}
		action, _ := GetPrompt("Select actions with task: "+taskId, actions)
		switch action {
		case 0:
			handlers.ViewTask(redis, taskId, taskType)
		case 1:
			handlers.EditTask(redis, taskId, taskType)
		case 2:
			handlers.DeleteTask(redis, taskId, DisplayTasks)
		case 3:
			DisplayTasks(redis)
		}
	}

}

func GetPrompt(label string, items []string) (int, string) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	index, name, err := prompt.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) {
			return -1, ""
		}
		log.Fatalf("Prompt failed %v\n", err)
	}
	return index, name
}
