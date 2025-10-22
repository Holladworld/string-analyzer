package handlers

import (
    "net/http"
    "strconv"
    "strings"
    "regexp"
    "reflect"
    "github.com/holladworld/string-analyzer/models"
    "github.com/holladworld/string-analyzer/services"
    "github.com/holladworld/string-analyzer/database"
    "github.com/gin-gonic/gin"
)

func PostStringHandler(c *gin.Context) {
    var request struct {
        Value interface{} `json:"value" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body or missing 'value' field"})
        return
    }
    
    if reflect.TypeOf(request.Value).Kind() != reflect.String {
        c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid data type for 'value' (must be string)"})
        return
    }
    
    stringValue := request.Value.(string)
    
    exists, err := database.StringExists(stringValue)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
    if exists {
        c.JSON(http.StatusConflict, gin.H{"error": "String already exists in the system"})
        return
    }
    
    result := services.AnalyzeString(stringValue)
    
    err = database.StoreString(result)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store string"})
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "id":    result.ID,
        "value": result.Value,
        "properties": gin.H{
            "length": result.Length,
            "is_palindrome": result.IsPalindrome,
            "unique_characters": result.UniqueCharacters,
            "word_count": result.WordCount,
            "sha256_hash": result.SHA256Hash,
            "character_frequency_map": result.CharacterFrequencyMap,
        },
        "created_at": result.CreatedAt,
    })
}

func GetStringHandler(c *gin.Context) {
    requestedValue := c.Param("string_value")
    
    result, exists, err := database.GetString(requestedValue)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
    if !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "String does not exist in the system"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "id":    result.ID,
        "value": result.Value,
        "properties": gin.H{
            "length": result.Length,
            "is_palindrome": result.IsPalindrome,
            "unique_characters": result.UniqueCharacters,
            "word_count": result.WordCount,
            "sha256_hash": result.SHA256Hash,
            "character_frequency_map": result.CharacterFrequencyMap,
        },
        "created_at": result.CreatedAt,
    })
}

func GetAllStringsHandler(c *gin.Context) {
    // Get query parameters
    isPalindromeStr := c.Query("is_palindrome")
    minLengthStr := c.Query("min_length")
    maxLengthStr := c.Query("max_length")
    wordCountStr := c.Query("word_count")
    containsChar := c.Query("contains_character")
    
    // Validate query parameters
    if isPalindromeStr != "" {
        if _, err := strconv.ParseBool(isPalindromeStr); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value for 'is_palindrome' (must be true or false)"})
            return
        }
    }
    
    if minLengthStr != "" {
        if _, err := strconv.Atoi(minLengthStr); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value for 'min_length' (must be integer)"})
            return
        }
    }
    
    if maxLengthStr != "" {
        if _, err := strconv.Atoi(maxLengthStr); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value for 'max_length' (must be integer)"})
            return
        }
    }
    
    if wordCountStr != "" {
        if _, err := strconv.Atoi(wordCountStr); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value for 'word_count' (must be integer)"})
            return
        }
    }
    
    if containsChar != "" && len(containsChar) != 1 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value for 'contains_character' (must be single character)"})
        return
    }
    
    allStrings, err := database.GetAllStrings()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
    
    filteredStrings := make([]models.AnalysisResult, 0)
    
    // Apply filters
    for _, str := range allStrings {
        if !applyFilters(str, isPalindromeStr, minLengthStr, maxLengthStr, wordCountStr, containsChar) {
            continue
        }
        filteredStrings = append(filteredStrings, str)
    }
    
    // Build response
    filtersApplied := gin.H{}
    if isPalindromeStr != "" { filtersApplied["is_palindrome"] = isPalindromeStr }
    if minLengthStr != "" { filtersApplied["min_length"] = minLengthStr }
    if maxLengthStr != "" { filtersApplied["max_length"] = maxLengthStr }
    if wordCountStr != "" { filtersApplied["word_count"] = wordCountStr }
    if containsChar != "" { filtersApplied["contains_character"] = containsChar }
    
    c.JSON(http.StatusOK, gin.H{
        "data": filteredStrings,
        "count": len(filteredStrings),
        "filters_applied": filtersApplied,
    })
}

func DeleteStringHandler(c *gin.Context) {
    requestedValue := c.Param("string_value")
    
    deleted, err := database.DeleteString(requestedValue)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
    if !deleted {
        c.JSON(http.StatusNotFound, gin.H{"error": "String does not exist in the system"})
        return
    }
    
    c.Status(http.StatusNoContent)
}

func NaturalLanguageFilterHandler(c *gin.Context) {
    query := c.Query("query")
    if query == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'query' is required"})
        return
    }
    
    // Parse natural language
    filters := parseNaturalLanguage(query)
    
    // Check for conflicting filters
    if filters.MinLength != nil && filters.MaxLength != nil && *filters.MinLength > *filters.MaxLength {
        c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Conflicting filters: min_length cannot be greater than max_length"})
        return
    }
    
    allStrings, err := database.GetAllStrings()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
    
    filteredStrings := make([]models.AnalysisResult, 0)
    
    for _, str := range allStrings {
        if !applyFilters(str, 
            boolToString(filters.IsPalindrome),
            intToString(filters.MinLength),
            intToString(filters.MaxLength),
            intToString(filters.WordCount),
            filters.ContainsCharacter) {
            continue
        }
        filteredStrings = append(filteredStrings, str)
    }
    
    c.JSON(http.StatusOK, gin.H{
        "data": filteredStrings,
        "count": len(filteredStrings),
        "interpreted_query": gin.H{
            "original": query,
            "parsed_filters": filters,
        },
    })
}

// Helper functions
func applyFilters(str models.AnalysisResult, isPalindromeStr, minLengthStr, maxLengthStr, wordCountStr, containsChar string) bool {
    // Palindrome filter
    if isPalindromeStr != "" {
        wanted, _ := strconv.ParseBool(isPalindromeStr)
        if str.IsPalindrome != wanted {
            return false
        }
    }
    
    // Min length filter
    if minLengthStr != "" {
        minLen, _ := strconv.Atoi(minLengthStr)
        if str.Length < minLen {
            return false
        }
    }
    
    // Max length filter
    if maxLengthStr != "" {
        maxLen, _ := strconv.Atoi(maxLengthStr)
        if str.Length > maxLen {
            return false
        }
    }
    
    // Word count filter
    if wordCountStr != "" {
        wordCount, _ := strconv.Atoi(wordCountStr)
        if str.WordCount != wordCount {
            return false
        }
    }
    
    // Contains character filter
    if containsChar != "" && len(containsChar) > 0 {
        if !strings.Contains(str.Value, containsChar) {
            return false
        }
    }
    
    return true
}

type NaturalLanguageFilters struct {
    IsPalindrome      *bool
    MinLength         *int
    MaxLength         *int
    WordCount         *int
    ContainsCharacter string
}

func parseNaturalLanguage(query string) NaturalLanguageFilters {
    filters := NaturalLanguageFilters{}
    query = strings.ToLower(query)
    
    // Palindrome detection
    if strings.Contains(query, "palindrom") {
        trueVal := true
        filters.IsPalindrome = &trueVal
    }
    
    // Word count detection
    if strings.Contains(query, "single word") || strings.Contains(query, "one word") {
        one := 1
        filters.WordCount = &one
    } else if strings.Contains(query, "two words") || strings.Contains(query, "2 words") {
        two := 2
        filters.WordCount = &two
    } else if strings.Contains(query, "three words") || strings.Contains(query, "3 words") {
        three := 3
        filters.WordCount = &three
    }
    
    // Length detection with number extraction
    if strings.Contains(query, "longer than") || strings.Contains(query, "more than") || strings.Contains(query, "at least") {
        re := regexp.MustCompile(`(\d+)`)
        matches := re.FindStringSubmatch(query)
        if len(matches) > 0 {
            if num, err := strconv.Atoi(matches[1]); err == nil {
                filters.MinLength = &num
            }
        }
    }
    
    if strings.Contains(query, "shorter than") || strings.Contains(query, "less than") || strings.Contains(query, "at most") {
        re := regexp.MustCompile(`(\d+)`)
        matches := re.FindStringSubmatch(query)
        if len(matches) > 0 {
            if num, err := strconv.Atoi(matches[1]); err == nil {
                filters.MaxLength = &num
            }
        }
    }
    
    // Enhanced character detection with vowels
    if strings.Contains(query, "contain") || strings.Contains(query, "with") || strings.Contains(query, "having") {
        // Vowel detection
        if strings.Contains(query, "first vowel") || strings.Contains(query, "vowel a") {
            filters.ContainsCharacter = "a"
        } else if strings.Contains(query, "vowel e") {
            filters.ContainsCharacter = "e"
        } else if strings.Contains(query, "vowel i") {
            filters.ContainsCharacter = "i"
        } else if strings.Contains(query, "vowel o") {
            filters.ContainsCharacter = "o"
        } else if strings.Contains(query, "vowel u") {
            filters.ContainsCharacter = "u"
        } else if strings.Contains(query, "letter z") || strings.Contains(query, "z") {
            filters.ContainsCharacter = "z"
        } else if strings.Contains(query, "letter a") || strings.Contains(query, " a ") {
            filters.ContainsCharacter = "a"
        } else if strings.Contains(query, "letter e") || strings.Contains(query, " e ") {
            filters.ContainsCharacter = "e"
        }
    }
    
    return filters
}

func boolToString(b *bool) string {
    if b == nil {
        return ""
    }
    return strconv.FormatBool(*b)
}

func intToString(i *int) string {
    if i == nil {
        return ""
    }
    return strconv.Itoa(*i)
}
