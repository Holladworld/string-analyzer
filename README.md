# String Analyzer API

A RESTful API service that analyzes strings and stores their computed properties. Built for HNG Backend Stage 1.

## 🚀 Live Demo

**API Base URL:** https://string-api-1761135512.fly.dev

## Features

- **POST /strings** - Analyze and store string properties
- **GET /strings/{value}** - Retrieve specific string analysis  
- **GET /strings** - Get all strings with advanced filtering
- **GET /strings/filter-by-natural-language** - Natural language query support
- **DELETE /strings/{value}** - Remove strings from storage

## Quick Test

```bash
curl -X POST https://string-api-1761135512.fly.dev/strings \
  -H "Content-Type: application/json" \
  -d '{"value": "hello world"}'cat > README.md << 'EOF'
# String Analyzer API

A RESTful API service that analyzes strings and stores their computed properties. Built for HNG Backend Stage 1.

## 🚀 Live Demo

**API Base URL:** https://string-api-1761135512.fly.dev

## Features

- **POST /strings** - Analyze and store string properties
- **GET /strings/{value}** - Retrieve specific string analysis  
- **GET /strings** - Get all strings with advanced filtering
- **GET /strings/filter-by-natural-language** - Natural language query support
- **DELETE /strings/{value}** - Remove strings from storage

## Quick Test

```bash
curl -X POST https://string-api-1761135512.fly.dev/strings \
  -H "Content-Type: application/json" \
  -d '{"value": "hello world"}'
🛠️ Tech Stack
Backend: Go 1.22

Framework: Gin Web Framework

Database: SQLite

Deployment: Fly.io
#API Endpoints
POST /strings
Analyze a string and store its properties.

Request:

json
{
  "value": "hello world"
}
Response (201 Created):

json
{
  "id": "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
  "value": "hello world",
  "properties": {
    "length": 11,
    "is_palindrome": false,
    "unique_characters": 8,
    "word_count": 2,
    "sha256_hash": "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
    "character_frequency_map": {
      "h": 1, "e": 1, "l": 3, "o": 2, " ": 1, "w": 1, "r": 1, "d": 1
    }
  },
  "created_at": "2024-01-21T10:00:00Z"
}
GET /strings/{string_value}
Retrieve analysis for a specific string.

GET /strings
Get all strings with optional filtering.

Query Parameters:

is_palindrome (boolean)

min_length (integer)

max_length (integer)

word_count (integer)

contains_character (string, single character)

GET /strings/filter-by-natural-language
Natural language query support.

Examples:

"all single word palindromic strings"

"strings longer than 10 characters"

"strings containing the letter z"

DELETE /strings/{string_value}
Remove a string from storage.

GitHub Repository
https://github.com/holladworld/string-analyzer

📄 License
This project is for HNG Internship submission.
