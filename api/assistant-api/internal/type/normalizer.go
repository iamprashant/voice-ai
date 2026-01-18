// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_type

import (
	"context"
	"strings"

	internal_normalizers "github.com/rapidaai/api/assistant-api/internal/normalizers"
	"github.com/rapidaai/pkg/commons"
)

// =============================================================================
// Text Normalizer Interface
// =============================================================================

// TextNormalizer defines the contract for provider-specific text preprocessing.
// Each TTS provider can implement this interface to handle its own normalization
// requirements (SSML, phonemes, abbreviations, etc.)
type TextNormalizer interface {
	// Normalize transforms text for optimal TTS output.
	Normalize(ctx context.Context, text string) string
}

// =============================================================================
// SSML Format Types
// =============================================================================

// SSMLFormat represents different SSML dialects supported by TTS providers.
type SSMLFormat string

const (
	SSMLFormatNone   SSMLFormat = "none"   // No SSML support (Deepgram, Cartesia)
	SSMLFormatW3C    SSMLFormat = "w3c"    // Standard W3C SSML
	SSMLFormatAzure  SSMLFormat = "azure"  // Azure Cognitive Services SSML
	SSMLFormatGoogle SSMLFormat = "google" // Google Cloud TTS SSML
	SSMLFormatAmazon SSMLFormat = "amazon" // Amazon Polly SSML
)

type NormalizerConfig struct {

	//
	Abbrieviations []string

	//
	Conjunctions []string

	//
	PauseDurationMs uint64
}

func DefaultNormalizerConfig() NormalizerConfig {
	return NormalizerConfig{
		Abbrieviations:  []string{},
		Conjunctions:    []string{},
		PauseDurationMs: 240,
	}
}

func BuildNormalizerPipeline(logger commons.Logger, names []string) []internal_normalizers.Normalizer {
	normalizers := make([]internal_normalizers.Normalizer, 0, len(names))

	for _, name := range names {
		name = strings.TrimSpace(strings.ToLower(name))
		var normalizer internal_normalizers.Normalizer

		switch name {
		case "url":
			normalizer = internal_normalizers.NewUrlNormalizer(logger)
		case "currency":
			normalizer = internal_normalizers.NewCurrencyNormalizer(logger)
		case "date":
			normalizer = internal_normalizers.NewDateNormalizer(logger)
		case "time":
			normalizer = internal_normalizers.NewTimeNormalizer(logger)
		case "number", "number-to-word":
			normalizer = internal_normalizers.NewNumberToWordNormalizer(logger)
		case "symbol":
			normalizer = internal_normalizers.NewSymbolNormalizer(logger)
		case "general-abbreviation", "general":
			normalizer = internal_normalizers.NewGeneralAbbreviationNormalizer(logger)
		case "role-abbreviation", "role":
			normalizer = internal_normalizers.NewRoleAbbreviationNormalizer(logger)
		case "tech-abbreviation", "tech":
			normalizer = internal_normalizers.NewTechAbbreviationNormalizer(logger)
		case "address":
			normalizer = internal_normalizers.NewAddressNormalizer(logger)
		default:
			logger.Warnf("normalizer: unknown normalizer '%s', skipping", name)
			continue
		}
		normalizers = append(normalizers, normalizer)
	}
	return normalizers
}
