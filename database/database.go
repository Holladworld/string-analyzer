package database

import (
    "database/sql"
    "fmt"
    "github.com/holladworld/string-analyzer/models"
    "encoding/json"
    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() error {
    var err error
    DB, err = sql.Open("sqlite3", "./string_analyzer.db")
    if err != nil {
        return err
    }
    
    fmt.Println("Connected to SQLite database")
    return createTable()
}

func createTable() error {
    query := `
    CREATE TABLE IF NOT EXISTS analyzed_strings (
        id TEXT PRIMARY KEY,
        value TEXT UNIQUE NOT NULL,
        length INTEGER NOT NULL,
        is_palindrome BOOLEAN NOT NULL,
        unique_characters INTEGER NOT NULL,
        word_count INTEGER NOT NULL,
        sha256_hash TEXT NOT NULL,
        character_frequency_map TEXT NOT NULL,
        created_at TEXT NOT NULL
    )
    `
    _, err := DB.Exec(query)
    return err
}

func StoreString(result models.AnalysisResult) error {
    freqMapJSON, err := json.Marshal(result.CharacterFrequencyMap)
    if err != nil {
        return err
    }
    
    query := `
    INSERT INTO analyzed_strings 
    (id, value, length, is_palindrome, unique_characters, word_count, sha256_hash, character_frequency_map, created_at)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
    _, err = DB.Exec(query, 
        result.ID, result.Value, result.Length, result.IsPalindrome,
        result.UniqueCharacters, result.WordCount, result.SHA256Hash,
        string(freqMapJSON), result.CreatedAt)
    return err
}

func GetString(value string) (models.AnalysisResult, bool, error) {
    var result models.AnalysisResult
    var freqMapJSON string
    
    query := "SELECT id, value, length, is_palindrome, unique_characters, word_count, sha256_hash, character_frequency_map, created_at FROM analyzed_strings WHERE value = ?"
    err := DB.QueryRow(query, value).Scan(
        &result.ID, &result.Value, &result.Length, &result.IsPalindrome,
        &result.UniqueCharacters, &result.WordCount, &result.SHA256Hash,
        &freqMapJSON, &result.CreatedAt)
    
    if err == sql.ErrNoRows {
        return result, false, nil
    }
    if err != nil {
        return result, false, err
    }
    
    // Parse JSON frequency map
    err = json.Unmarshal([]byte(freqMapJSON), &result.CharacterFrequencyMap)
    if err != nil {
        return result, false, err
    }
    
    return result, true, nil
}

func GetAllStrings() ([]models.AnalysisResult, error) {
    query := "SELECT id, value, length, is_palindrome, unique_characters, word_count, sha256_hash, character_frequency_map, created_at FROM analyzed_strings"
    rows, err := DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var results []models.AnalysisResult
    for rows.Next() {
        var result models.AnalysisResult
        var freqMapJSON string
        
        err := rows.Scan(
            &result.ID, &result.Value, &result.Length, &result.IsPalindrome,
            &result.UniqueCharacters, &result.WordCount, &result.SHA256Hash,
            &freqMapJSON, &result.CreatedAt)
        if err != nil {
            return nil, err
        }
        
        err = json.Unmarshal([]byte(freqMapJSON), &result.CharacterFrequencyMap)
        if err != nil {
            return nil, err
        }
        
        results = append(results, result)
    }
    
    return results, nil
}

func DeleteString(value string) (bool, error) {
    query := "DELETE FROM analyzed_strings WHERE value = ?"
    result, err := DB.Exec(query, value)
    if err != nil {
        return false, err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return false, err
    }
    
    return rowsAffected > 0, nil
}

func StringExists(value string) (bool, error) {
    var exists bool
    query := "SELECT EXISTS(SELECT 1 FROM analyzed_strings WHERE value = ?)"
    err := DB.QueryRow(query, value).Scan(&exists)
    return exists, err
}
