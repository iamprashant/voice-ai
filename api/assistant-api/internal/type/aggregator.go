// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_type

import (
	"context"
)

// LLMTextAggregator defines the contract for components that transform
// streamed or batched text inputs into aggregated sentence outputs.
//
// Implementations are expected to:
//   - Accept inputs via Aggregate
//   - Emit results asynchronously on the Result channel
//   - Release resources and stop processing on Close
type LLMTextAggregator interface {
	// Aggregate consumes an aggregator input (such as an LLMStreamChunk
	// or Finalize signal). Implementations should respect context
	// cancellation and deadlines.
	Aggregate(ctx context.Context, in ...LLMPacket) error

	// Result returns a read-only channel on which aggregated outputs
	// are delivered.
	Result() <-chan Packet

	// Close terminates the aggregator, releases resources,
	// and closes the Result channel.
	Close() error
}
