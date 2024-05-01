package party

import (
	"fmt"
	"sync"
)

// PlanParty uses goroutines to complete tasks for a party you are hosting
//
// Example:
//
//	PlanParty()
//
// This function identifies tasks for the party, such as sending invites, cleaning, buying refreshments,
// creating a playlist, and buying a new outfit. It launches a goroutine for each task and waits for
// all tasks to complete before displaying the results and starting the party.
func PlanParty() {
	// Identify tasks for the party
	tasks := []string{"Send invites", "Cleaning", "Buy refreshments", "Create playlist", "Buy new outfit"}

	// Create a channel to communicate tasks completed.
	taskResults := make(chan string)

	// Create a WaitGroup to await all goroutines completing.
	var waitGroup sync.WaitGroup

	// Start a goroutine for each individual task
	for _, task := range tasks {
		waitGroup.Add(1)
		go func(task string) {
			defer waitGroup.Done()
			taskResults <- doTask(task)
		}(task)
	}

	// Await the completion of all the party tasks and ensure values are received by channel
	go func() {
		waitGroup.Wait()
		close(taskResults)
	}()

	// Display the results of the tasks
	for result := range taskResults {
		fmt.Println(result)
	}

	fmt.Println("\nLet's get this party started...")
}

func doTask(task string) string {
	var result string

	if task == "Buy new outfit" {
		result = task + ": There's no budget for this!"
	} else {
		result = task + " done!"
	}

	return result
}
