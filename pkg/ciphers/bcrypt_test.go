// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package ciphers

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash_BasicFunctionality(t *testing.T) {
	input := "test string"
	result := Hash(input)

	// Should return a non-empty string
	assert.NotEmpty(t, result)

	// Should be 32 characters (MD5 hex length)
	assert.Len(t, result, 32)

	// Should only contain hexadecimal characters
	assert.Regexp(t, regexp.MustCompile("^[a-f0-9]+$"), result)
}

func TestHash_Deterministic(t *testing.T) {
	input := "deterministic test"
	result1 := Hash(input)
	result2 := Hash(input)
	result3 := Hash(input)

	// Should always return the same result for the same input
	assert.Equal(t, result1, result2)
	assert.Equal(t, result2, result3)
	assert.Equal(t, result1, result3)
}

func TestHash_EmptyString(t *testing.T) {
	result := Hash("")

	// Should handle empty string
	assert.NotEmpty(t, result)
	assert.Len(t, result, 32)
	assert.Regexp(t, regexp.MustCompile("^[a-f0-9]+$"), result)

	// MD5 of empty string should be d41d8cd98f00b204e9800998ecf8427e
	expected := "d41d8cd98f00b204e9800998ecf8427e"
	assert.Equal(t, expected, result)
}

func TestHash_SpecialCharacters(t *testing.T) {
	inputs := []string{
		"Hello 🌍 World!",
		"Special chars: @#$%^&*()",
		"Unicode: 你好世界",
		"Newlines\nand\ttabs",
		"Quotes: 'single' and \"double\"",
	}

	for _, input := range inputs {
		result := Hash(input)
		assert.Len(t, result, 32)
		assert.Regexp(t, regexp.MustCompile("^[a-f0-9]+$"), result)
	}
}

func TestHash_LongString(t *testing.T) {
	// Create a very long string (1MB)
	longString := strings.Repeat("This is a test string that will be repeated many times. ", 20000)
	assert.Greater(t, len(longString), 1000000) // > 1MB

	result := Hash(longString)
	assert.Len(t, result, 32)
	assert.Regexp(t, regexp.MustCompile("^[a-f0-9]+$"), result)
}

func TestHash_DifferentInputs(t *testing.T) {
	inputs := []string{
		"hello",
		"world",
		"Hello",
		"hello ",
		" hello",
		"123",
		"abc",
	}

	results := make(map[string]bool)
	for _, input := range inputs {
		result := Hash(input)
		// All results should be different (collision test)
		assert.False(t, results[result], "Hash collision detected for input: %s", input)
		results[result] = true
		assert.Len(t, result, 32)
	}
}

func TestHash_CompareWithStandardMD5(t *testing.T) {
	input := "test input for comparison"
	result := Hash(input)

	// Compute MD5 using standard library
	expected := md5.Sum([]byte(input))
	expectedHex := hex.EncodeToString(expected[:])

	assert.Equal(t, expectedHex, result)
}

func TestRandomHash_BasicFunctionality(t *testing.T) {
	prefix := "test"
	result := RandomHash(prefix)

	// Should return a non-empty string
	assert.NotEmpty(t, result)

	// Should be 32 characters (MD5 hex length)
	assert.Len(t, result, 32)

	// Should only contain hexadecimal characters
	assert.Regexp(t, regexp.MustCompile("^[a-f0-9]+$"), result)
}

func TestRandomHash_NonDeterministic(t *testing.T) {
	prefix := "random_test"
	result1 := RandomHash(prefix)
	result2 := RandomHash(prefix)
	result3 := RandomHash(prefix)

	// Should return different results each time (due to UUID)
	assert.NotEqual(t, result1, result2)
	assert.NotEqual(t, result2, result3)
	assert.NotEqual(t, result1, result3)
}

func TestRandomHash_EmptyPrefix(t *testing.T) {
	result := RandomHash("")

	// Should handle empty prefix
	assert.NotEmpty(t, result)
	assert.Len(t, result, 32)
	assert.Regexp(t, regexp.MustCompile("^[a-f0-9]+$"), result)
}

func TestRandomHash_SpecialCharactersPrefix(t *testing.T) {
	prefixes := []string{
		"prefix_with_spaces",
		"prefix-with-dashes",
		"prefix.with.dots",
		"prefix@domain.com",
		"prefix#tag",
		"prefix$var",
	}

	for _, prefix := range prefixes {
		result := RandomHash(prefix)
		assert.Len(t, result, 32)
		assert.Regexp(t, regexp.MustCompile("^[a-f0-9]+$"), result)
	}
}

func TestRandomHash_LongPrefix(t *testing.T) {
	longPrefix := strings.Repeat("very_long_prefix_", 100)
	result := RandomHash(longPrefix)

	assert.Len(t, result, 32)
	assert.Regexp(t, regexp.MustCompile("^[a-f0-9]+$"), result)
}

func TestRandomHash_UnicodePrefix(t *testing.T) {
	prefixes := []string{
		"префикс",
		"前缀",
		"プレフィックス",
		"πρόθεμα",
	}

	for _, prefix := range prefixes {
		result := RandomHash(prefix)
		assert.Len(t, result, 32)
		assert.Regexp(t, regexp.MustCompile("^[a-f0-9]+$"), result)
	}
}

func TestToken_BasicFunctionality(t *testing.T) {
	prefix := "token_test"
	result := Token(prefix)

	// Should return a non-empty string
	assert.NotEmpty(t, result)

	// Should be 64 characters (SHA256 hex length)
	assert.Len(t, result, 64)

	// Should only contain hexadecimal characters
	assert.Regexp(t, regexp.MustCompile("^[a-f0-9]+$"), result)
}

func TestToken_NonDeterministic(t *testing.T) {
	prefix := "token_random"
	result1 := Token(prefix)
	result2 := Token(prefix)
	result3 := Token(prefix)

	// Should return different results each time
	assert.NotEqual(t, result1, result2)
	assert.NotEqual(t, result2, result3)
	assert.NotEqual(t, result1, result3)
}

func TestToken_EmptyPrefix(t *testing.T) {
	result := Token("")

	// Should handle empty prefix
	assert.NotEmpty(t, result)
	assert.Len(t, result, 64)
	assert.Regexp(t, regexp.MustCompile("^[a-f0-9]+$"), result)
}

func TestToken_SpecialCharactersPrefix(t *testing.T) {
	prefixes := []string{
		"token_with_spaces",
		"token-with-dashes",
		"token.with.dots",
		"token@domain.com",
		"token#tag",
		"token$var",
	}

	for _, prefix := range prefixes {
		result := Token(prefix)
		assert.Len(t, result, 64)
		assert.Regexp(t, regexp.MustCompile("^[a-f0-9]+$"), result)
	}
}

func TestToken_LongPrefix(t *testing.T) {
	longPrefix := strings.Repeat("very_long_token_prefix_", 100)
	result := Token(longPrefix)

	assert.Len(t, result, 64)
	assert.Regexp(t, regexp.MustCompile("^[a-f0-9]+$"), result)
}

func TestToken_UnicodePrefix(t *testing.T) {
	prefixes := []string{
		"токен",
		"令牌",
		"トークン",
		"κουπόνι",
	}

	for _, prefix := range prefixes {
		result := Token(prefix)
		assert.Len(t, result, 64)
		assert.Regexp(t, regexp.MustCompile("^[a-f0-9]+$"), result)
	}
}

func TestToken_CompareWithStandardSHA256(t *testing.T) {
	prefix := "comparison_test"
	token := Token(prefix)

	// Token should be 64 characters (SHA256 hex)
	assert.Len(t, token, 64)

	// Should be valid hex
	assert.Regexp(t, regexp.MustCompile("^[a-f0-9]+$"), token)

	// Since Token internally calls RandomHash, we can verify it's SHA256 of some 32-char hex string
	// But we can't predict the exact value due to UUID randomness
	// Instead, let's verify the token is different from the prefix hash
	prefixHash := Hash(prefix)
	assert.NotEqual(t, prefixHash, token)
	assert.Len(t, prefixHash, 32) // MD5
	assert.Len(t, token, 64)      // SHA256
}

func TestToken_DifferentFromRandomHash(t *testing.T) {
	prefix := "comparison"

	randomHash := RandomHash(prefix)
	token := Token(prefix)

	// Token and RandomHash should be different (different hash algorithms)
	assert.NotEqual(t, randomHash, token)
	assert.Len(t, randomHash, 32) // MD5
	assert.Len(t, token, 64)      // SHA256
}

func TestHash_Vs_RandomHash_Different(t *testing.T) {
	input := "test_input"

	hashResult := Hash(input)
	randomHashResult := RandomHash(input)

	// Hash and RandomHash should produce different results
	// Hash is deterministic MD5 of input
	// RandomHash is MD5 of "input_uuid"
	assert.NotEqual(t, hashResult, randomHashResult)
}

func TestIntegration_Hash_RandomHash_Token(t *testing.T) {
	prefix := "integration_test"

	// Generate a specific random hash
	randomHash := RandomHash(prefix)

	// Create token manually using that random hash
	expectedToken := sha256.Sum256([]byte(randomHash))
	expectedTokenHex := hex.EncodeToString(expectedToken[:])

	// Verify that if we manually compute SHA256 of a random hash, we get a valid token format
	assert.Len(t, expectedTokenHex, 64)
	assert.Regexp(t, regexp.MustCompile("^[a-f0-9]+$"), expectedTokenHex)

	// Verify that Token() produces a valid token (though we can't predict the exact value)
	actualToken := Token(prefix)
	assert.Len(t, actualToken, 64)
	assert.Regexp(t, regexp.MustCompile("^[a-f0-9]+$"), actualToken)

	// Verify that Token and RandomHash produce different outputs
	assert.NotEqual(t, randomHash, actualToken)
	assert.Len(t, randomHash, 32)  // MD5
	assert.Len(t, actualToken, 64) // SHA256
}

func TestHash_CaseSensitivity(t *testing.T) {
	input1 := "Hello"
	input2 := "hello"
	input3 := "HELLO"

	result1 := Hash(input1)
	result2 := Hash(input2)
	result3 := Hash(input3)

	// MD5 is case-sensitive, so these should all be different
	assert.NotEqual(t, result1, result2)
	assert.NotEqual(t, result2, result3)
	assert.NotEqual(t, result1, result3)
}

func TestRandomHash_PrefixInclusion(t *testing.T) {
	prefix := "unique_prefix_123"
	randomHash := RandomHash(prefix)

	// While we can't predict the exact hash, we can verify it's different from just hashing the prefix
	prefixOnlyHash := Hash(prefix)
	assert.NotEqual(t, prefixOnlyHash, randomHash)
}

func TestToken_Uniqueness(t *testing.T) {
	prefix := "uniqueness_test"

	// Generate multiple tokens and ensure they're all unique
	tokens := make(map[string]bool)
	for i := 0; i < 100; i++ {
		token := Token(prefix)
		assert.False(t, tokens[token], "Token collision detected")
		tokens[token] = true
	}

	assert.Len(t, tokens, 100)
}

func TestRandomHash_Uniqueness(t *testing.T) {
	prefix := "uniqueness_test"

	// Generate multiple random hashes and ensure they're all unique
	hashes := make(map[string]bool)
	for i := 0; i < 100; i++ {
		hash := RandomHash(prefix)
		assert.False(t, hashes[hash], "RandomHash collision detected")
		hashes[hash] = true
	}

	assert.Len(t, hashes, 100)
}

func TestHash_KnownValues(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"", "d41d8cd98f00b204e9800998ecf8427e"},
		{"a", "0cc175b9c0f1b6a831c399e269772661"},
		{"hello", "5d41402abc4b2a76b9719d911017c592"},
		{"hello world", "5eb63bbbe01eeed093cb22bb8f5acdc3"},
		{"The quick brown fox jumps over the lazy dog", "9e107d9d372bb6826bd81d3542a419d6"},
	}

	for _, tc := range testCases {
		result := Hash(tc.input)
		assert.Equal(t, tc.expected, result, "Hash mismatch for input: %s", tc.input)
	}
}

func TestHash_Performance(t *testing.T) {
	// Test with large input to ensure no performance issues
	largeInput := strings.Repeat("This is a performance test string. ", 100000) // ~3MB

	result := Hash(largeInput)
	assert.Len(t, result, 32)
	assert.Regexp(t, regexp.MustCompile("^[a-f0-9]+$"), result)
}

func TestToken_FormatValidation(t *testing.T) {
	prefix := "format_test"

	// Test that token format is always valid hex
	for i := 0; i < 10; i++ {
		token := Token(prefix)
		assert.Len(t, token, 64)
		assert.Regexp(t, regexp.MustCompile("^[a-f0-9]+$"), token)

		// Should not contain uppercase letters (hex.EncodeToString produces lowercase)
		assert.NotRegexp(t, regexp.MustCompile("[A-F]"), token)
	}
}

func TestRandomHash_FormatValidation(t *testing.T) {
	prefix := "format_test"

	// Test that random hash format is always valid hex
	for i := 0; i < 10; i++ {
		hash := RandomHash(prefix)
		assert.Len(t, hash, 32)
		assert.Regexp(t, regexp.MustCompile("^[a-f0-9]+$"), hash)

		// Should not contain uppercase letters
		assert.NotRegexp(t, regexp.MustCompile("[A-F]"), hash)
	}
}
