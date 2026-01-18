// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_transformer_azure

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
// Azure Text Normalizer
// =============================================================================

// azureNormalizer handles Azure Cognitive Services TTS text preprocessing.
// Azure supports full SSML with mstts extensions for expressive speech.
type azureNormalizer struct {
	logger    commons.Logger
	config    internal_type.NormalizerConfig
	voiceName string
	language  string

	// normalizer pipeline
	normalizers []internal_normalizers.Normalizer

	// conjunction handling
	conjunctionPattern *regexp.Regexp
}

// NewAzureNormalizer creates an Azure-specific text normalizer.
func NewAzureNormalizer(logger commons.Logger, opts utils.Option) internal_type.TextNormalizer {
	cfg := internal_type.DefaultNormalizerConfig()

	// Get voice name and language
	voiceName, _ := opts.GetString("speaker.voice.name")
	language, _ := opts.GetString("speaker.language")
	if language == "" {
		language = "en-US"
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

	return &azureNormalizer{
		logger:             logger,
		config:             cfg,
		voiceName:          voiceName,
		language:           language,
		normalizers:        normalizers,
		conjunctionPattern: conjunctionPattern,
	}
}

// Normalize applies Azure-specific text transformations.
func (n *azureNormalizer) Normalize(ctx context.Context, text string) string {
	if text == "" {
		return text
	}

	// Clean markdown first
	text = n.removeMarkdown(text)

	// Apply normalizer pipeline
	for _, normalizer := range n.normalizers {
		text = normalizer.Normalize(text)
	}

	// Escape XML special characters for SSML safety (Azure uses SSML)
	text = n.escapeXML(text)

	// Insert breaks after conjunction boundaries
	if n.conjunctionPattern != nil && n.config.PauseDurationMs > 0 {
		text = n.insertConjunctionBreaks(text)
	}

	return n.normalizeWhitespace(text)
}

// =============================================================================
// Private Helpers
// =============================================================================

func (n *azureNormalizer) removeMarkdown(input string) string {
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

func (n *azureNormalizer) escapeXML(text string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&apos;",
	)
	return replacer.Replace(text)
}

func (n *azureNormalizer) insertConjunctionBreaks(text string) string {
	breakTag := fmt.Sprintf(`<break time="%dms"/>`, n.config.PauseDurationMs)
	return n.conjunctionPattern.ReplaceAllStringFunc(text, func(match string) string {
		return match + breakTag
	})
}

func (n *azureNormalizer) normalizeWhitespace(text string) string {
	re := regexp.MustCompile(`\s+`)
	result := re.ReplaceAllString(text, " ")
	return strings.TrimSpace(result)
}

// =============================================================================
// Azure SSML Helpers
// =============================================================================

func (n *azureNormalizer) WrapWithSSML(text string) string {
	return fmt.Sprintf(
		`<speak version="1.0" xmlns="http://www.w3.org/2001/10/synthesis" xmlns:mstts="https://www.w3.org/2001/mstts" xml:lang="%s"><voice name="%s">%s</voice></speak>`,
		n.language, n.voiceName, text,
	)
}

func (n *azureNormalizer) AddBreak(durationMs int) string {
	return fmt.Sprintf(`<break time="%dms"/>`, durationMs)
}

func (n *azureNormalizer) AddProsody(text string, rate, pitch, volume string) string {
	attrs := ""
	if rate != "" {
		attrs += fmt.Sprintf(` rate="%s"`, rate)
	}
	if pitch != "" {
		attrs += fmt.Sprintf(` pitch="%s"`, pitch)
	}
	if volume != "" {
		attrs += fmt.Sprintf(` volume="%s"`, volume)
	}
	if attrs == "" {
		return text
	}
	return fmt.Sprintf(`<prosody%s>%s</prosody>`, attrs, text)
}

func (n *azureNormalizer) AddEmphasis(text, level string) string {
	return fmt.Sprintf(`<emphasis level="%s">%s</emphasis>`, level, text)
}

func (n *azureNormalizer) AddExpressAs(text, style string) string {
	return fmt.Sprintf(`<mstts:express-as style="%s">%s</mstts:express-as>`, style, text)
}

func (n *azureNormalizer) SayAs(text, interpretAs, format string) string {
	if format != "" {
		return fmt.Sprintf(`<say-as interpret-as="%s" format="%s">%s</say-as>`, interpretAs, format, text)
	}
	return fmt.Sprintf(`<say-as interpret-as="%s">%s</say-as>`, interpretAs, text)
}
