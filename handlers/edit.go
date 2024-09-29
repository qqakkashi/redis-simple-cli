package handlers

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
)

func EditTask(client *redis.Client, taskId string, taskType string) {
	switch taskType {
	case "hash":
		EditHashTask(client, taskId)
	case "set":
		EditSetTask(client, taskId)
	case "zset":
		EditZSetTask(client, taskId)
	case "string":
		EditStringTask(client, taskId)
	}
}

func EditHashTask(client *redis.Client, taskID string) {
	for {
		currentData, err := client.HGetAll(ctx, taskID).Result()
		if err != nil {
			log.Fatalf("Error getting current task data: %s\n", err)
			return
		}

		if len(currentData) == 0 {
			log.Fatalf("Task not found.")
			return
		}

		fields := make([]string, 0, len(currentData))
		for field := range currentData {
			fields = append(fields, field)
		}

		promptField := promptui.Select{
			Label: "Select a field to edit",
			Items: fields,
		}

		_, selectedField, err := promptField.Run()
		if err != nil {
			log.Fatalf("Prompt failed: %s\n", err)
			return
		}

		currentValue := currentData[selectedField]
		log.Printf("Current value for %s: %s\n", selectedField, currentValue)

		promptValue := promptui.Prompt{
			Label:   "Enter new value",
			Default: currentValue,
		}

		newValue, err := promptValue.Run()
		if err != nil {
			log.Fatalf("Prompt failed: %s", err)
			return
		}

		promptConfirm := promptui.Select{
			Label: fmt.Sprintf("Apply changes for field %s to %s?", selectedField, newValue),
			Items: []string{"Yes", "No"},
		}

		_, confirmation, err := promptConfirm.Run()
		if err != nil {
			log.Println("Prompt failed:", err)
			return
		}

		if confirmation == "Yes" {
			err = client.HSet(ctx, taskID, selectedField, newValue).Err()
			if err != nil {
				log.Fatalf("Error updating task: %s \n", err)
				return
			}
			log.Printf("Task %s updated successfully. Field %s set to %s\n", taskID, selectedField, newValue)
		} else {
			log.Println("Changes were not applied.")
		}

		if ContinuePrompt("Do you want to edit another field?") == false {
			break
		}
	}
}

func EditSetTask(client *redis.Client, taskID string) {
	for {
		currentMembers, err := client.SMembers(ctx, taskID).Result()
		if err != nil {
			log.Fatalf("Error getting current set members: %s\n", err)
			return
		}

		if len(currentMembers) == 0 {
			log.Fatalf("Task not found or empty set.")
			return
		}

		promptField := promptui.Select{
			Label: "Select a member to edit",
			Items: currentMembers,
		}

		_, selectedMember, err := promptField.Run()
		if err != nil {
			log.Fatalf("Prompt failed: %s\n", err)
			return
		}

		fmt.Printf("Current member: %s\n", selectedMember)

		promptValue := promptui.Prompt{
			Label:   "Enter new value (or leave blank to remove)",
			Default: selectedMember,
		}

		newValue, err := promptValue.Run()
		if err != nil {
			log.Println("Prompt failed:", err)
			return
		}

		if newValue == "" {
			err = client.SRem(ctx, taskID, selectedMember).Err()
			if err != nil {
				log.Fatalf("Error removing member: %s\n", err)
			} else {
				log.Printf("Member %s removed successfully from task %s\n", selectedMember, taskID)
			}
		} else {
			err = client.SRem(ctx, taskID, selectedMember).Err()
			if err != nil {
				log.Fatalf("Error removing old member: %s\n", err)
			}
			err = client.SAdd(ctx, taskID, newValue).Err()
			if err != nil {
				log.Fatalf("Error adding new member: %s\n", err)
			} else {
				log.Printf("Member updated successfully. Old: %s, New: %s\n", selectedMember, newValue)
			}
		}

		if ContinuePrompt("Do you want to edit another member?") == false {
			break
		}

	}
}

func EditZSetTask(client *redis.Client, taskID string) {
	for {
		currentMembers, err := client.ZRange(ctx, taskID, 0, -1).Result()
		if err != nil {
			log.Fatalf("Error getting current zset members: %s\n", err)
			return
		}

		if len(currentMembers) == 0 {
			log.Fatalf("Task not found or empty zset.")
			return
		}

		promptField := promptui.Select{
			Label: "Select a member to edit",
			Items: currentMembers,
		}

		_, selectedMember, err := promptField.Run()
		if err != nil {
			log.Fatalf("Prompt failed: %s\n", err)
			return
		}

		fmt.Printf("Current member: %s\n", selectedMember)

		promptValue := promptui.Prompt{
			Label:   "Enter new score",
			Default: "0",
		}

		newScoreStr, err := promptValue.Run()
		if err != nil {
			log.Println("Prompt failed:", err)
			return
		}

		newScore, err := strconv.ParseFloat(newScoreStr, 64)
		if err != nil {
			log.Println("Invalid score input:", err)
			return
		}
		z := redis.Z{Score: newScore, Member: selectedMember}

		err = client.ZAdd(ctx, taskID, z).Err()
		if err != nil {
			log.Fatalf("Error updating zset member: %s\n", err)
		} else {
			log.Printf("Member %s updated successfully with new score %f\n", selectedMember, newScore)
		}

		if ContinuePrompt("Do you want to edit another member?") == false {
			break
		}
	}
}

func EditStringTask(client *redis.Client, taskID string) {
	for {
		currentValue, err := client.Get(ctx, taskID).Result()
		if err != nil {
			if err == redis.Nil {
				log.Fatalf("Task not found.")
			}
			log.Fatalf("Error getting current task data: %s\n", err)
			return
		}

		fmt.Printf("Current value for task %s: %s\n", taskID, currentValue)

		promptValue := promptui.Prompt{
			Label:   "Enter new value",
			Default: currentValue,
		}

		newValue, err := promptValue.Run()
		if err != nil {
			log.Println("Prompt failed:", err)
			return
		}

		promptConfirm := promptui.Select{
			Label: fmt.Sprintf("Apply changes for task %s to %s?", taskID, newValue),
			Items: []string{"Yes", "No"},
		}

		_, confirmation, err := promptConfirm.Run()
		if err != nil {
			log.Println("Prompt failed:", err)
			return
		}

		if confirmation == "Yes" {
			err = client.Set(ctx, taskID, newValue, 0).Err()
			if err != nil {
				log.Fatalf("Error updating task: %s\n", err)
				return
			}
			log.Printf("Task %s updated successfully. New value set to %s\n", taskID, newValue)
		} else {
			log.Println("Changes were not applied.")
		}

		if ContinuePrompt("Do you want to edit another task?") == false {
			break
		}

	}
}

func ContinuePrompt(label string) bool {
	promptContinue := promptui.Select{
		Label: label,
		Items: []string{"Yes", "No"},
	}

	_, continueEditing, err := promptContinue.Run()
	if err != nil || continueEditing == "No" {
		return false
	}
	return true
}
