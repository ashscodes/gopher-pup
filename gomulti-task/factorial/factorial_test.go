package factorial

import (
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"gomulti-task/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock big.Int to isolate the behaviour of calculatePartialFactorial
type MockBigInt struct {
	mock.Mock
}

func (m *MockBigInt) Mul(x, y *big.Int) *big.Int {
	args := m.Called(x, y)
	return args.Get(0).(*big.Int)
}

// TestCalculateFactorial tests CalculateFactorial
func TestCalculateFactorial(t *testing.T) {
	testCases := []struct {
		ExpectedError  error
		ExpectedResult *big.Int
		Name           string
		Number         int
		TotalRoutines  int
	}{
		{nil, big.NewInt(120), "When_number_is_5_and_total_goroutines_is_2", 5, 2},
		{nil, big.NewInt(3628800), "When_number_is_10_and_total_goroutines_is_3", 10, 3},
		{nil, big.NewInt(1), "When_number_is_0_and_total_goroutines_is_1", 0, 1},
		{errors.New("input number must be non-negative"), nil, "When_number_is_-5_and_total_goroutines_is_2", -5, 2},
		{errors.New("totalRoutines must be greater than zero"), nil, "When_number_is_5_and_total_goroutines_is_0", 5, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := CalculateFactorial(tc.Number, tc.TotalRoutines)
			assert.Equal(t, tc.ExpectedResult, result)
			assert.Equal(t, tc.ExpectedError, err)
		})
	}
}

// TestCalculateFactorialTermination tests that the CalculateFactorial function terminates correctly.
func TestCalculateFactorialTermination(t *testing.T) {
	// Use a channel to signal completion
	done := make(chan struct{})

	go func() {
		defer close(done)
		CalculateFactorial(10, 5)
	}()

	// Wait for the function to finish or timeout after a specified duration
	timeout := 5 * time.Second
	select {
	case <-done:
		// Function completed successfully
	case <-time.After(timeout):
		t.Error("CalculateFactorial did not terminate within the expected time")
	}
}

func BenchmarkCalculateFactorial(b *testing.B) {
	numbers := []int{1, 5, 10, 15, 20}
	routines := []int{1, 2, 4, 8, 16}

	// disable standard output.
	helpers.DisableStdOut()

	for _, number := range numbers {
		for _, routine := range routines {
			b.Run(fmt.Sprintf("When_Number_is_%d_and_total_goroutines_is_%d", number, routine), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					CalculateFactorial(number, routine)
				}
			})
		}
	}

	// Re-enable standard output.
	helpers.EnableStdOut()
}
