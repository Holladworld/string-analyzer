package models

type AnalysisResult struct {
    ID                    string         `json:"id"`
    Value                 string         `json:"value"`
    Length                int            `json:"length"`
    IsPalindrome          bool           `json:"is_palindrome"`
    UniqueCharacters      int            `json:"unique_characters"`
    WordCount             int            `json:"word_count"`
    SHA256Hash            string         `json:"sha256_hash"`
    CharacterFrequencyMap map[string]int `json:"character_frequency_map"` // Changed to string keys
    CreatedAt             string         `json:"created_at"`
}
