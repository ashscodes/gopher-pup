package factorial

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
)

// CalculateFactorial uses goroutines to calcualte the factorial value of a given number.
//
// Parameters:
//
//	number: Integer for which to calculate the factorial value.
//	totalRoutines: Number of goroutines launched.
//
// Returns:
//   - The factorial of the input number as a *big.Int.
//   - An error, if any, encountered during the calculation.
//
// Example:
//
//	// Calculate the factorial of 10 using 5 goroutines.
//	result, err := CalculateFactorial(10, 5)
//	if err != nil {
//	    fmt.Println("Error:", err)
//	} else {
//	    fmt.Println("Factorial result:", result)
//	}
func CalculateFactorial(number, totalRoutines int) (*big.Int, error) {
	if number < 0 {
		return nil, errors.New("input number must be non-negative")
	}

	if totalRoutines <= 0 {
		return nil, errors.New("totalRoutines must be greater than zero")
	}

	fmt.Printf("Calculating factorial of %d using %d goroutines.\n", number, totalRoutines)

	// Create channel to receive factorialResults from the goroutines.
	factorialResults := make(chan *big.Int, totalRoutines)

	// Create a WaitGroup to await all goroutines completing.
	var waitGroup sync.WaitGroup
	waitGroup.Add(totalRoutines)

	// Calculate partial factorial results concurrently.
	for i := 0; i < totalRoutines; i++ {
		go func(i int) {
			defer waitGroup.Done()

			start := 1 + (i * (number / totalRoutines))
			end := (i + 1) * (number / totalRoutines)
			if i == totalRoutines-1 {
				end = number
			}
			factorialResults <- calculatePartialFactorial(start, end)
		}(i)
	}

	// Await the completion of all the partial factorial results and ensure values are received by channel
	go func() {
		waitGroup.Wait()
		close(factorialResults)
	}()

	// Aggregate partial results
	finalResult := big.NewInt(1)
	for partialResult := range factorialResults {
		finalResult.Mul(finalResult, partialResult)
	}

	fmt.Printf("Factorial of %d is %d\n", number, finalResult)

	return finalResult, nil
}

func calculatePartialFactorial(start, end int) *big.Int {
	result := big.NewInt(1)
	for i := start; i <= end; i++ {
		result.Mul(result, big.NewInt(int64(i)))
	}

	return result
}
