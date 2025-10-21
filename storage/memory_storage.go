package storage

import (
    "sync"
    "string-analyzer/models"
)

var (
    stringsMap = make(map[string]models.AnalysisResult)
    mutex      = &sync.RWMutex{}
)

// StoreString saves an analyzed string
func StoreString(result models.AnalysisResult) {
    mutex.Lock()
    defer mutex.Unlock()
    stringsMap[result.Value] = result
}

// GetString retrieves a string by its value
func GetString(value string) (models.AnalysisResult, bool) {
    mutex.RLock()
    defer mutex.RUnlock()
    result, exists := stringsMap[value]
    return result, exists
}

// GetAllStrings returns all stored strings
func GetAllStrings() []models.AnalysisResult {
    mutex.RLock()
    defer mutex.RUnlock()
    
    results := make([]models.AnalysisResult, 0, len(stringsMap))
    for _, result := range stringsMap {
        results = append(results, result)
    }
    return results
}

// DeleteString removes a string by its value
func DeleteString(value string) bool {
    mutex.Lock()
    defer mutex.Unlock()
    
    if _, exists := stringsMap[value]; exists {
        delete(stringsMap, value)
        return true
    }
    return false
}

// StringExists checks if a string already exists
func StringExists(value string) bool {
    mutex.RLock()
    defer mutex.RUnlock()
    _, exists := stringsMap[value]
    return exists
}