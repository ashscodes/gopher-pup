package main

import (
	"gomulti-task/factorial"
	"gomulti-task/party"
)

// Concurrency examples
func main() {
	factorial.CalculateFactorial(10, 5)
	party.PlanParty()
}
