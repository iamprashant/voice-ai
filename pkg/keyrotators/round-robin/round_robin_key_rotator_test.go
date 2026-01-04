// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package roundrobin_keyrotators

import (
	"sync"
	"testing"

	"github.com/rapidaai/pkg/commons"
	"github.com/stretchr/testify/assert"
)

func TestNewRoundRobinKeyRotator(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	keys := []string{"key1", "key2", "key3"}

	rotator := NewRoundRobinKeyRotator(logger, keys...)

	assert.NotNil(t, rotator)
	assert.Equal(t, keys, rotator.GetAll())
}

func TestNewRoundRobinKeyRotator_EmptyKeys(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()

	rotator := NewRoundRobinKeyRotator(logger)

	assert.NotNil(t, rotator)
	assert.Empty(t, rotator.GetAll())
}

func TestRoundRobinKeyRotator_Next_EmptyKeys(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	rotator := NewRoundRobinKeyRotator(logger)

	result := rotator.Next()

	assert.Empty(t, result)
}

func TestRoundRobinKeyRotator_Next_SingleKey(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	rotator := NewRoundRobinKeyRotator(logger, "single-key")

	result := rotator.Next()

	assert.Equal(t, "single-key", result)
}

func TestRoundRobinKeyRotator_Next_MultipleKeys(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	keys := []string{"key1", "key2", "key3"}
	rotator := NewRoundRobinKeyRotator(logger, keys...)

	// Test sequential rotation
	assert.Equal(t, "key2", rotator.Next()) // counter starts at 0, +1 = 1 -> key2
	assert.Equal(t, "key3", rotator.Next()) // +1 = 2 -> key3
	assert.Equal(t, "key1", rotator.Next()) // +1 = 3 % 3 = 0 -> key1
	assert.Equal(t, "key2", rotator.Next()) // +1 = 1 -> key2
}

func TestRoundRobinKeyRotator_Add(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	rotator := NewRoundRobinKeyRotator(logger, "key1")

	err := rotator.Add("key2")
	assert.NoError(t, err)
	err = rotator.Add("key3")
	assert.NoError(t, err)

	assert.Equal(t, []string{"key1", "key2", "key3"}, rotator.GetAll())
}

func TestRoundRobinKeyRotator_Remove_Existing(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	rotator := NewRoundRobinKeyRotator(logger, "key1", "key2", "key3")

	err := rotator.Remove("key2")
	assert.NoError(t, err)

	assert.Equal(t, []string{"key1", "key3"}, rotator.GetAll())
}

func TestRoundRobinKeyRotator_Remove_FirstElement(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	rotator := NewRoundRobinKeyRotator(logger, "key1", "key2", "key3")

	err := rotator.Remove("key1")
	assert.NoError(t, err)

	assert.Equal(t, []string{"key2", "key3"}, rotator.GetAll())
}

func TestRoundRobinKeyRotator_Remove_LastElement(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	rotator := NewRoundRobinKeyRotator(logger, "key1", "key2", "key3")

	err := rotator.Remove("key3")
	assert.NoError(t, err)

	assert.Equal(t, []string{"key1", "key2"}, rotator.GetAll())
}

func TestRoundRobinKeyRotator_Remove_NonExisting(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	originalKeys := []string{"key1", "key2", "key3"}
	rotator := NewRoundRobinKeyRotator(logger, originalKeys...)

	err := rotator.Remove("non-existing")
	assert.NoError(t, err)

	assert.Equal(t, originalKeys, rotator.GetAll())
}

func TestRoundRobinKeyRotator_Remove_EmptySlice(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	rotator := NewRoundRobinKeyRotator(logger)

	err := rotator.Remove("any-key")
	assert.NoError(t, err)

	assert.Empty(t, rotator.GetAll())
}

func TestRoundRobinKeyRotator_Remove_DuplicateKeys(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	rotator := NewRoundRobinKeyRotator(logger, "key1", "key2", "key1", "key3")

	err := rotator.Remove("key1")
	assert.NoError(t, err)

	// Should remove only the first occurrence
	assert.Equal(t, []string{"key2", "key1", "key3"}, rotator.GetAll())
}

func TestRoundRobinKeyRotator_GetAll(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	originalKeys := []string{"key1", "key2", "key3"}
	rotator := NewRoundRobinKeyRotator(logger, originalKeys...)

	result := rotator.GetAll()

	assert.Equal(t, originalKeys, result)
	// Ensure it's a copy, not the original slice
	result[0] = "modified"
	assert.Equal(t, []string{"key1", "key2", "key3"}, rotator.GetAll())
}

func TestRoundRobinKeyRotator_GetAll_Empty(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	rotator := NewRoundRobinKeyRotator(logger)

	result := rotator.GetAll()

	assert.Empty(t, result)
}

func TestRoundRobinKeyRotator_ThreadSafety_Next(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	keys := []string{"key1", "key2", "key3", "key4", "key5"}
	rotator := NewRoundRobinKeyRotator(logger, keys...)

	var wg sync.WaitGroup
	results := make(chan string, 100)
	numGoroutines := 10
	numCallsPerGoroutine := 10

	// Start multiple goroutines calling Next concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numCallsPerGoroutine; j++ {
				result := rotator.Next()
				results <- result
			}
		}()
	}

	wg.Wait()
	close(results)

	// Collect all results
	var allResults []string
	for result := range results {
		allResults = append(allResults, result)
	}

	assert.Len(t, allResults, numGoroutines*numCallsPerGoroutine)

	// Verify all results are valid keys
	for _, result := range allResults {
		assert.Contains(t, keys, result)
	}
}
