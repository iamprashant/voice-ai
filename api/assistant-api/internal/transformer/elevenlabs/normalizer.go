// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_transformer_elevenlabs

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	internal_normalizers "github.com/rapidaai/api/assistant-api/internal/normalizers"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/utils"
)

// =============================================================================
// ElevenLabs Text Normalizer
// =============================================================================

// elevenlabsNormalizer handles ElevenLabs TTS text preprocessing.
// ElevenLabs supports LIMITED SSML: only <break> and <phoneme> tags.
// Break time is specified in SECONDS (e.g., time="0.5s" not "500ms").
type elevenlabsNormalizer struct {
	logger   commons.Logger
	config   internal_type.NormalizerConfig
	language string

	// normalizer pipeline
	normalizers []internal_normalizers.Normalizer

	// conjunction handling
	conjunctionPattern *regexp.Regexp
}

// NewElevenLabsNormalizer creates an ElevenLabs-specific text normalizer.
func NewElevenLabsNormalizer(logger commons.Logger, opts utils.Option) internal_type.TextNormalizer {
	cfg := internal_type.DefaultNormalizerConfig()

	language, _ := opts.GetString("speaker.language")
	if language == "" {
		language = "en"
	}

	// Parse conjunction boundaries from options
	var conjunctionPattern *regexp.Regexp
	if conjunctionBoundaries, err := opts.GetString("speaker.conjunction.boundaries"); err == nil && conjunctionBoundaries != "" {
		cfg.Conjunctions = strings.Split(conjunctionBoundaries, commons.SEPARATOR)

		escaped := make([]string, len(cfg.Conjunctions))
		for i, c := range cfg.Conjunctions {
			escaped[i] = regexp.QuoteMeta(strings.TrimSpace(c))
		}
		pattern := `(` + strings.Join(escaped, "|") + `)`
		conjunctionPattern = regexp.MustCompile(pattern)
	}

	// Parse conjunction break duration
	if conjunctionBreak, err := opts.GetUint64("speaker.conjunction.break"); err == nil {
		cfg.PauseDurationMs = conjunctionBreak
	}

	// Build normalizer pipeline based on speaker.pronunciation.dictionaries
	var normalizers []internal_normalizers.Normalizer
	if dictionaries, err := opts.GetString("speaker.pronunciation.dictionaries"); err == nil && dictionaries != "" {
		normalizerNames := strings.Split(dictionaries, commons.SEPARATOR)
		normalizers = internal_type.BuildNormalizerPipeline(logger, normalizerNames)
	}

	return &elevenlabsNormalizer{
		logger:             logger,
		config:             cfg,
		language:           language,
		normalizers:        normalizers,
		conjunctionPattern: conjunctionPattern,
	}
}

// Normalize applies ElevenLabs-specific text transformations.
// ElevenLabs supports only <break> and <phoneme> SSML tags.
func (n *elevenlabsNormalizer) Normalize(ctx context.Context, text string) string {
	if text == "" {
		return text
	}

	// Clean markdown first
	text = n.removeMarkdown(text)

	// Apply normalizer pipeline
	for _, normalizer := range n.normalizers {
		text = normalizer.Normalize(text)
	}

	// ElevenLabs supports limited SSML, so we escape XML characters
	// except where we insert our own SSML tags
	text = n.escapeXML(text)

	// Insert breaks after conjunction boundaries (ElevenLabs uses seconds)
	if n.conjunctionPattern != nil && n.config.PauseDurationMs > 0 {
		text = n.insertConjunctionBreaks(text)
	}

	return n.normalizeWhitespace(text)
}

// =============================================================================
// Private Helpers
// =============================================================================

func (n *elevenlabsNormalizer) removeMarkdown(input string) string {
	re := regexp.MustCompile(`(?m)^#{1,6}\s*`)
	output := re.ReplaceAllString(input, "")

	re = regexp.MustCompile(`\*{1,2}([^*]+?)\*{1,2}|_{1,2}([^_]+?)_{1,2}`)
	output = re.ReplaceAllString(output, "$1$2")

	re = regexp.MustCompile("`([^`]+)`")
	output = re.ReplaceAllString(output, "$1")

	re = regexp.MustCompile("(?s)```[^`]*```")
	output = re.ReplaceAllString(output, "")

	re = regexp.MustCompile(`(?m)^>\s?`)
	output = re.ReplaceAllString(output, "")

	re = regexp.MustCompile(`\[(.*?)\]\(.*?\)`)
	output = re.ReplaceAllString(output, "$1")

	re = regexp.MustCompile(`!\[(.*?)\]\(.*?\)`)
	output = re.ReplaceAllString(output, "$1")

	re = regexp.MustCompile(`(?m)^(-{3,}|\*{3,}|_{3,})$`)
	output = re.ReplaceAllString(output, "")

	re = regexp.MustCompile(`[*_]+`)
	output = re.ReplaceAllString(output, "")

	return output
}

// escapeXML escapes XML special characters for limited SSML safety.
func (n *elevenlabsNormalizer) escapeXML(text string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
	)
	return replacer.Replace(text)
}

// insertConjunctionBreaks adds breaks after conjunctions.
// ElevenLabs uses seconds format (e.g., "0.5s" instead of "500ms").
func (n *elevenlabsNormalizer) insertConjunctionBreaks(text string) string {
	// Convert milliseconds to seconds for ElevenLabs
	seconds := float64(n.config.PauseDurationMs) / 1000.0
	breakTag := fmt.Sprintf(`<break time="%.2fs"/>`, seconds)

	return n.conjunctionPattern.ReplaceAllStringFunc(text, func(match string) string {
		return match + breakTag
	})
}

func (n *elevenlabsNormalizer) normalizeWhitespace(text string) string {
	re := regexp.MustCompile(`\s+`)
	result := re.ReplaceAllString(text, " ")
	return strings.TrimSpace(result)
}

// =============================================================================
// ElevenLabs SSML Helpers (Limited Support)
// =============================================================================

// AddBreak adds a pause. ElevenLabs uses seconds format.
func (n *elevenlabsNormalizer) AddBreak(durationMs int) string {
	seconds := float64(durationMs) / 1000.0
	return fmt.Sprintf(`<break time="%.2fs"/>`, seconds)
}

// AddPhoneme wraps text with phoneme pronunciation.
// ElevenLabs supports IPA phoneme alphabet.
func (n *elevenlabsNormalizer) AddPhoneme(text, ipa string) string {
	return fmt.Sprintf(`<phoneme alphabet="ipa" ph="%s">%s</phoneme>`, ipa, text)
}
