// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_default_aggregator

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/utils"
)

// textAggregator implements sentence-level aggregation for streaming text.
type textAggregator struct {
	logger commons.Logger
	ctx    context.Context

	result chan internal_type.Packet

	mu sync.Mutex // Changed from RWMutex - we rarely have read-only operations

	// buffering state
	buffer         strings.Builder
	currentContext string
	boundaryRegex  *regexp.Regexp
	hasBoundaries  bool

	// optimization: reusable slice to reduce allocations
	toEmitBuffer []internal_type.Packet
}

// NewDefaultLLMTextAggregator creates a new sentence aggregator with the given logger and options.
//
// The options parameter should contain "speaker.sentence.boundaries" which is a
// comma-separated list of sentence delimiters (e.g., ".,?!").
//
// Returns an error if the boundary regex compilation fails.
//
// Example:
//
//	aggregator, err := NewDefaultLLMTextAggregator(logger, options)
//	if err != nil {
//		return err
//	}
func NewDefaultLLMTextAggregator(context context.Context, logger commons.Logger, options utils.Option) (internal_type.LLMTextAggregator, error) {
	st := &textAggregator{
		ctx:          context,
		logger:       logger,
		result:       make(chan internal_type.Packet, 32), // Increased buffer for better throughput
		toEmitBuffer: make([]internal_type.Packet, 0, 8),  // Pre-allocate slice
	}
	if err := st.initializeBoundaries(options, logger); err != nil {
		return nil, err
	}
	return st, nil
}

// initializeBoundaries sets up the boundary regex from configuration options.
// It safely handles missing or invalid boundary configurations.
func (st *textAggregator) initializeBoundaries(options utils.Option, logger commons.Logger) error {
	boundariesRaw, err := options.GetString("speaker.sentence.boundaries")
	if err != nil || boundariesRaw == "" {
		return nil
	}

	boundaries := strings.Split(boundariesRaw, commons.SEPARATOR)
	validBoundaries := filterBoundaries(boundaries)

	if len(validBoundaries) == 0 {
		return nil
	}
	st.hasBoundaries = true
	parts := make([]string, 0, len(validBoundaries))
	for _, b := range validBoundaries {
		parts = append(parts, regexp.QuoteMeta(b))
	}

	pattern := fmt.Sprintf(`(%s)\s*`, strings.Join(parts, "|"))
	regex, err := regexp.Compile(pattern)
	if err != nil {
		logger.Errorf("Invalid boundary regex: %v", err)
		st.hasBoundaries = false
		return nil
	}

	st.boundaryRegex = regex
	logger.Debugf("Custom sentence boundaries: %v", validBoundaries)
	return nil
}

// filterBoundaries removes empty and whitespace-only boundaries.
func filterBoundaries(boundaries []string) []string {
	var valid []string
	for _, b := range boundaries {
		if b = strings.TrimSpace(b); b != "" {
			valid = append(valid, b)
		}
	}
	return valid
}

// Aggregate processes a text segment and emits complete sentences.
//
// It handles context switching by flushing any pending text in the previous context.
// When sentence.IsComplete is true, any remaining buffered text is emitted as a complete sentence.
//
// The method respects context cancellation and returns immediately if ctx is cancelled.
//
// Example:
//
//	err := aggregator.Aggregate(ctx, Text{
//		ContextId:  "speaker1",
//		Text:   "Hello world.",
//		IsComplete: true,
//	})
func (st *textAggregator) Aggregate(ctx context.Context, sentences ...internal_type.LLMPacket) error {
	// Process all sentences with lock held once, then emit without lock
	st.mu.Lock()
	st.toEmitBuffer = st.toEmitBuffer[:0] // Reset reusable buffer

	for _, sentence := range sentences {
		st.extractAndQueueLocked(sentence)
	}

	// Copy results to avoid holding lock during channel operations
	toEmit := make([]internal_type.Packet, len(st.toEmitBuffer))
	copy(toEmit, st.toEmitBuffer)
	st.mu.Unlock()

	// Send queued results while respecting context cancellation
	for _, s := range toEmit {
		select {
		case st.result <- s:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

// extractAndQueueLocked extracts complete sentences from the input.
// MUST be called with lock held. Appends results to st.toEmitBuffer.
func (st *textAggregator) extractAndQueueLocked(sentence internal_type.LLMPacket) {
	switch input := sentence.(type) {
	case internal_type.LLMResponseDeltaPacket:
		// Handle context switch - just clean buffer, do NOT emit
		if input.ContextID != st.currentContext && st.currentContext != "" {
			st.buffer.Reset()
		}

		st.currentContext = input.ContextID
		st.buffer.WriteString(input.Text)

		// Extract sentences at boundaries
		if st.hasBoundaries {
			st.extractTextsByBoundaryLocked(input.ContextID)
		}
	case internal_type.LLMResponseDonePacket:
		// Flush remaining buffer
		if st.buffer.Len() > 0 {
			content := st.buffer.String()
			if content != "" {
				st.toEmitBuffer = append(st.toEmitBuffer, internal_type.LLMResponseDeltaPacket{
					ContextID: input.ContextID,
					Text:      content,
				})
			}
			st.buffer.Reset()
		}
		st.toEmitBuffer = append(st.toEmitBuffer, input)
	default:
		st.logger.Warnf("Unsupported tokenizer input type: %T", sentence)
	}
}

// extractTextsByBoundaryLocked extracts all sentences that end at boundaries.
// MUST be called with lock held. Appends to st.toEmitBuffer.
func (st *textAggregator) extractTextsByBoundaryLocked(contextId string) {
	text := st.buffer.String()

	// Find all boundaries at once instead of iterating
	matches := st.boundaryRegex.FindAllStringIndex(text, -1)
	if len(matches) == 0 {
		return
	}

	// Extract complete sentences up to the last boundary
	lastBoundary := matches[len(matches)-1][1]
	if lastBoundary > 0 {
		completeText := strings.TrimSpace(text[:lastBoundary])
		if completeText != "" {
			st.toEmitBuffer = append(st.toEmitBuffer, internal_type.LLMResponseDeltaPacket{
				ContextID: contextId,
				Text:      completeText,
			})
		}

		// Keep only the remaining text after last boundary
		st.buffer.Reset()
		if lastBoundary < len(text) {
			st.buffer.WriteString(text[lastBoundary:])
		}
	}
}

// Result returns a read-only channel that receives complete sentences.
//
// The channel is closed when Close() is called. Callers should use a range loop
// to safely consume sentences until the channel closes.
//
// Example:
//
//	for sentence := range aggregator.Result() {
//		fmt.Println(sentence.Text)
//	}
func (st *textAggregator) Result() <-chan internal_type.Packet {
	return st.result
}

// Close gracefully shuts down the aggregator and closes the result channel.
//
// After Close() is called, the aggregator cannot be reused. Any subsequent calls
// to Aggregate() will panic when trying to send on the closed channel.
//
// It is safe to call Close() multiple times; subsequent calls are no-ops.
func (st *textAggregator) Close() error {
	st.mu.Lock()
	defer st.mu.Unlock()

	// Clear state
	st.buffer.Reset()
	st.currentContext = ""

	// Close result channel
	if st.result != nil {
		close(st.result)
		st.result = nil
	}

	return nil
}

// String returns a string representation of the aggregator for debugging.
func (st *textAggregator) String() string {
	st.mu.Lock()
	defer st.mu.Unlock()

	return fmt.Sprintf("TextAggregator{context=%q, bufferLen=%d, hasBoundaries=%v}", st.currentContext, st.buffer.Len(), st.hasBoundaries)
}
