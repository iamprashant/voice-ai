// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package types

import (
	"strings"
)

// Emotion represents a standardized user emotion
type Emotion struct {
	Name  string // e.g., "anger", "positivity", "sadness"
	Level string // e.g., "low", "high", "moderate"
}

// GetEmotionFromString maps raw user emotion strings to structured Emotion
func GetEmotionByName(raw string) Emotion {
	rawLower := strings.ToLower(strings.TrimSpace(raw))
	switch rawLower {
	case "angry", "annoyed", "irritated", "mad":
		return Emotion{"anger", "high"}
	case "frustrated":
		return Emotion{"anger", "low"}
	case "sad", "worried", "disappointed", "upset":
		return Emotion{"sadness", "high"}
	case "confused", "skeptical":
		return Emotion{"curiosity", "low"}
	case "curious", "interested":
		return Emotion{"curiosity", ""}
	case "happy", "excited", "grateful", "joyful":
		return Emotion{"positivity", "high"}
	case "content", "relieved", "calm":
		return Emotion{"positivity", "low"}
	default:
		return Emotion{"neutral", ""}
	}
}
