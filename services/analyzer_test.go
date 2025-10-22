package services

import (
    "crypto/sha256"
    "encoding/hex"
    "testing"
)

// TestAnalyzeString tests the string analysis functionality
func TestAnalyzeString(t *testing.T) {
    input := "hello"
    result := AnalyzeString(input)
    
    // Test basic properties
    if result.Length != 5 {
        t.Errorf("Expected length 5, got %d", result.Length)
    }
    
    if result.WordCount != 1 {
        t.Errorf("Expected word count 1, got %d", result.WordCount)
    }
    
    if result.IsPalindrome != false {
        t.Errorf("Expected not palindrome, got %t", result.IsPalindrome)
    }
    
    if result.UniqueCharacters != 4 {
        t.Errorf("Expected 4 unique characters, got %d", result.UniqueCharacters)
    }
    
    // Test SHA256 hash
    expectedHash := sha256.Sum256([]byte(input))
    expectedHashStr := hex.EncodeToString(expectedHash[:])
    if result.SHA256Hash != expectedHashStr {
        t.Errorf("SHA256 hash mismatch")
    }
}

// TestPalindrome tests palindrome detection
func TestPalindrome(t *testing.T) {
    result := AnalyzeString("madam")
    if !result.IsPalindrome {
        t.Error("'madam' should be detected as palindrome")
    }
    
    result2 := AnalyzeString("hello")
    if result2.IsPalindrome {
        t.Error("'hello' should not be detected as palindrome")
    }
}

// TestCharacterFrequency tests character frequency counting
func TestCharacterFrequency(t *testing.T) {
    result := AnalyzeString("hello")
    
    expectedFreq := map[string]int{
        "h": 1, "e": 1, "l": 2, "o": 1,
    }
    
    for char, expectedCount := range expectedFreq {
        if result.CharacterFrequencyMap[char] != expectedCount {
            t.Errorf("Character frequency mismatch for '%s': got %d, want %d", char, result.CharacterFrequencyMap[char], expectedCount)
        }
    }
}
