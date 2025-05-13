package standalone

import (

	"math/rand"
	"sync"
	"testing"
	"time"
)

// This is a simplified version of the generateUniqueLuckyNumber function
// for testing the thread safety mechanism
var luckyNumberMutex sync.Mutex
var usedNumbers = make(map[int]bool)

func generateUniqueLuckyNumber() int {
	// Lock the mutex to ensure thread safety
	luckyNumberMutex.Lock()
	defer luckyNumberMutex.Unlock()

	// Generate a random number between 1 and 1,000,000
	const maxNumber = 1000000
	var number int
	
	// Try up to 10 times to generate a unique number
	for i := 0; i < 10; i++ {
		number = rand.Intn(maxNumber) + 1
		
		// Check if the number is already used
		if !usedNumbers[number] {
			// Mark it as used
			usedNumbers[number] = true
			return number
		}
	}
	
	// If we couldn't find a unique number, generate one using timestamp
	number = int(time.Now().UnixNano() % maxNumber)
	if number < 100000 {
		number += 100000 // Ensure 6 digits
	}
	
	usedNumbers[number] = true
	return number
}

func TestConcurrentLuckyNumberGeneration(t *testing.T) {
	// Number of concurrent users
	const numUsers = 100000
	
	// Use a deterministic random source for reproducible tests
	rand.Seed(42)
	
	// Channels for results
	results := make(chan int, numUsers)
	
	// Create a wait group to wait for all goroutines
	var wg sync.WaitGroup
	wg.Add(numUsers)
	
	// Start goroutines to simulate concurrent users
	for i := 0; i < numUsers; i++ {
		go func() {
			defer wg.Done()
			
			// Generate a lucky number
			results <- generateUniqueLuckyNumber()
		}()
	}
	
	// Close channel when all goroutines are done
	go func() {
		wg.Wait()
		close(results)
	}()
	
	// Collect and verify lucky numbers
	luckyNumbers := make(map[int]bool)
	duplicates := 0
	
	for num := range results {
		if luckyNumbers[num] {
			duplicates++
			t.Logf("Duplicate lucky number detected: %d", num)
		}
		luckyNumbers[num] = true
	}
	
	// Verify all numbers are unique
	if duplicates > 0 {
		t.Errorf("Found %d duplicate lucky numbers", duplicates)
	} else {
		t.Logf("Successfully generated %d unique lucky numbers", len(luckyNumbers))
	}
	
	// Verify all numbers are within expected range
	for num := range luckyNumbers {
		if num <= 0 || num > 1000000 {
			t.Errorf("Lucky number %d is outside valid range", num)
		}
	}
}
