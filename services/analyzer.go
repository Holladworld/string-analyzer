package services

import (
    "crypto/sha256"
    "encoding/hex"
    "strings"
    "time"
    "string-analyzer/models"
)

func AnalyzeString(input string) models.AnalysisResult {
    result := models.AnalysisResult{
        Value: input,
        CharacterFrequencyMap: make(map[string]int), // Changed to string
        CreatedAt: time.Now().UTC().Format(time.RFC3339),
    }
    
    // 1. Basic length
    result.Length = len(input)
    
    // 2. Word count (split by spaces)
    words := strings.Fields(input)
    result.WordCount = len(words)
    
    // 3. Palindrome check (case insensitive)
    cleanString := strings.ToLower(strings.ReplaceAll(input, " ", ""))
    isPalindrome := true
    for i := 0; i < len(cleanString)/2; i++ {
        if cleanString[i] != cleanString[len(cleanString)-1-i] {
            isPalindrome = false
            break
        }
    }
    result.IsPalindrome = isPalindrome
    
    // 4. Character frequency and unique count (FIXED)
    uniqueChars := make(map[rune]bool)
    for _, char := range input {
        charStr := string(char)
        result.CharacterFrequencyMap[charStr]++
        uniqueChars[char] = true
    }
    result.UniqueCharacters = len(uniqueChars)
    
    // 5. SHA256 hash
    hash := sha256.Sum256([]byte(input))
    result.SHA256Hash = hex.EncodeToString(hash[:])
    result.ID = result.SHA256Hash
    
    return result
}