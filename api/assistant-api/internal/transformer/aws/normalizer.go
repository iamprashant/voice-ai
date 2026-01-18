// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_transformer_aws

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
// AWS Polly Text Normalizer
// =============================================================================

// awsNormalizer handles AWS Polly TTS text preprocessing.
// AWS Polly supports full SSML with Amazon-specific extensions.
type awsNormalizer struct {
	logger commons.Logger
	config internal_type.NormalizerConfig

	// normalizer pipeline
	normalizers []internal_normalizers.Normalizer

	// conjunction handling
	conjunctionPattern *regexp.Regexp
}

// NewAWSNormalizer creates an AWS Polly-specific text normalizer.
func NewAWSNormalizer(logger commons.Logger, opts utils.Option) internal_type.TextNormalizer {
	cfg := internal_type.DefaultNormalizerConfig()

	// Parse conjunction boundaries from options
	var conjunctionPattern *regexp.Regexp
	if conjunctionBoundaries, err := opts.GetString("speaker.conjunction.boundaries"); err == nil && conjunctionBoundaries != "" {
		cfg.Conjunctions = strings.Split(conjunctionBoundaries, commons.SEPARATOR)

		// Build conjunction pattern
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

	return &awsNormalizer{
		logger:             logger,
		config:             cfg,
		normalizers:        normalizers,
		conjunctionPattern: conjunctionPattern,
	}
}

// buildNormalizerPipeline creates normalizers based on the provided names.
// Supported normalizer names: url, currency, date, time, number, symbol,
// general-abbreviation, role-abbreviation, tech-abbreviation, address

// Normalize applies AWS Polly-specific text transformations.
func (n *awsNormalizer) Normalize(ctx context.Context, text string) string {
	if text == "" {
		return text
	}

	// Clean markdown first (always applied)
	text = n.removeMarkdown(text)

	// Apply normalizer pipeline (only if configured)
	for _, normalizer := range n.normalizers {
		text = normalizer.Normalize(text)
	}

	// Escape XML special characters for SSML safety
	text = n.escapeXML(text)
	// Insert breaks after conjunction boundaries (only if configured)
	if n.conjunctionPattern != nil && n.config.PauseDurationMs > 0 {
		text = n.insertConjunctionBreaks(text)
	}
	return n.normalizeWhitespace(text)
}

// =============================================================================
// Private Helpers
// =============================================================================

// removeMarkdown strips markdown formatting from text.
func (n *awsNormalizer) removeMarkdown(input string) string {
	// Remove headers (#, ##, ### ...)
	re := regexp.MustCompile(`(?m)^#{1,6}\s*`)
	output := re.ReplaceAllString(input, "")

	// Remove bold/italic markers (*, **, _, __)
	re = regexp.MustCompile(`\*{1,2}([^*]+?)\*{1,2}|_{1,2}([^_]+?)_{1,2}`)
	output = re.ReplaceAllString(output, "$1$2")

	// Remove inline code/backticks
	re = regexp.MustCompile("`([^`]+)`")
	output = re.ReplaceAllString(output, "$1")

	// Remove code blocks
	re = regexp.MustCompile("(?s)```[^`]*```")
	output = re.ReplaceAllString(output, "")

	// Remove blockquotes (>)
	re = regexp.MustCompile(`(?m)^>\s?`)
	output = re.ReplaceAllString(output, "")

	// Remove links [text](url) → keep text
	re = regexp.MustCompile(`\[(.*?)\]\(.*?\)`)
	output = re.ReplaceAllString(output, "$1")

	// Remove images ![alt](url) → keep alt
	re = regexp.MustCompile(`!\[(.*?)\]\(.*?\)`)
	output = re.ReplaceAllString(output, "$1")

	// Remove horizontal rules (---, ***)
	re = regexp.MustCompile(`(?m)^(-{3,}|\*{3,}|_{3,})$`)
	output = re.ReplaceAllString(output, "")

	// Remove extra asterisks/underscores
	re = regexp.MustCompile(`[*_]+`)
	output = re.ReplaceAllString(output, "")

	return output
}

// escapeXML escapes XML special characters for SSML.
func (n *awsNormalizer) escapeXML(text string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&apos;",
	)
	return replacer.Replace(text)
}

// insertConjunctionBreaks inserts SSML breaks after conjunctions.
func (n *awsNormalizer) insertConjunctionBreaks(text string) string {
	breakTag := fmt.Sprintf(`<break time="%dms"/>`, n.config.PauseDurationMs)
	return n.conjunctionPattern.ReplaceAllStringFunc(text, func(match string) string {
		return match + breakTag
	})
}

// normalizeWhitespace collapses multiple spaces and trims.
func (n *awsNormalizer) normalizeWhitespace(text string) string {
	re := regexp.MustCompile(`\s+`)
	result := re.ReplaceAllString(text, " ")
	return strings.TrimSpace(result)
}

// =============================================================================
// AWS Polly SSML Helpers
// =============================================================================

// WrapWithSSML wraps text in AWS Polly SSML structure.
func (n *awsNormalizer) WrapWithSSML(text string) string {
	return fmt.Sprintf(`<speak>%s</speak>`, text)
}

// AddBreak creates an AWS Polly SSML break element.
func (n *awsNormalizer) AddBreak(durationMs int) string {
	return fmt.Sprintf(`<break time="%dms"/>`, durationMs)
}

// AddBreakStrength creates a break with strength attribute.
// strength: none, x-weak, weak, medium, strong, x-strong
func (n *awsNormalizer) AddBreakStrength(strength string) string {
	return fmt.Sprintf(`<break strength="%s"/>`, strength)
}

// AddProsody wraps text with prosody controls.
// rate: x-slow, slow, medium, fast, x-fast, or percentage
// pitch: x-low, low, medium, high, x-high, or percentage
// volume: silent, x-soft, soft, medium, loud, x-loud, or dB
func (n *awsNormalizer) AddProsody(text string, rate, pitch, volume string) string {
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

// AddEmphasis wraps text with emphasis.
// level: strong, moderate, reduced
func (n *awsNormalizer) AddEmphasis(text, level string) string {
	return fmt.Sprintf(`<emphasis level="%s">%s</emphasis>`, level, text)
}

// SayAs specifies how to interpret text.
// interpretAs: characters, spell-out, cardinal, ordinal, digits, fraction, unit, date, time, telephone, address, interjection, expletive
func (n *awsNormalizer) SayAs(text, interpretAs, format string) string {
	if format != "" {
		return fmt.Sprintf(`<say-as interpret-as="%s" format="%s">%s</say-as>`, interpretAs, format, text)
	}
	return fmt.Sprintf(`<say-as interpret-as="%s">%s</say-as>`, interpretAs, text)
}

// AddAmazonEffect adds Amazon-specific voice effects.
// effect: drc (Dynamic Range Compression), whispered
func (n *awsNormalizer) AddAmazonEffect(text, effect string) string {
	return fmt.Sprintf(`<amazon:effect name="%s">%s</amazon:effect>`, effect, text)
}

// AddWhisper creates whispering effect.
func (n *awsNormalizer) AddWhisper(text string) string {
	return n.AddAmazonEffect(text, "whispered")
}

// AddDomain specifies the speaking domain for neural voices.
// domain: conversational, news, long-form
func (n *awsNormalizer) AddDomain(text, domain string) string {
	return fmt.Sprintf(`<amazon:domain name="%s">%s</amazon:domain>`, domain, text)
}

// AddPhoneme specifies pronunciation using phonetic alphabet.
// alphabet: ipa, x-sampa
func (n *awsNormalizer) AddPhoneme(text, phoneme, alphabet string) string {
	return fmt.Sprintf(`<phoneme alphabet="%s" ph="%s">%s</phoneme>`, alphabet, phoneme, text)
}

// AddLang specifies language for a section.
func (n *awsNormalizer) AddLang(text, lang string) string {
	return fmt.Sprintf(`<lang xml:lang="%s">%s</lang>`, lang, text)
}
