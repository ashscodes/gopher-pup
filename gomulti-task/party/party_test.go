package party

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"gomulti-task/helpers"

	"github.com/stretchr/testify/assert"
)

// TestPlanParty tests the PlanParty function.
func TestPlanParty(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call function
	PlanParty()

	// Restore stdout
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = old

	expectedLines := []string{
		"Send invites done!",
		"Cleaning done!",
		"Buy refreshments done!",
		"Create playlist done!",
		"Buy new outfit: There's no budget for this!",
		"",
		"Let's get this party started...",
	}

	actualLines := strings.Split(string(out), "\n")

	for _, actualLine := range actualLines {
		assert.Contains(t, expectedLines, actualLine)
	}
}

// TestPlanPartyTermination tests that the PlanParty function terminates correctly.
func TestPlanPartyTermination(t *testing.T) {
	// Use a channel to signal completion
	done := make(chan struct{})

	go func() {
		defer close(done)
		PlanParty()
	}()

	// Wait for the function to finish or timeout after a specified duration
	timeout := 5 * time.Second
	select {
	case <-done:
		// Function completed successfully
	case <-time.After(timeout):
		t.Error("PlanParty did not terminate within the expected time")
	}
}

func BenchmarkPlanParty(b *testing.B) {
	// disable standard output.
	helpers.DisableStdOut()

	for i := 0; i < b.N; i++ {
		PlanParty()
	}

	// Re-enable standard output.
	helpers.EnableStdOut()
}
