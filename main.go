package main

import (
    "log"
    "os"
    "github.com/holladworld/string-analyzer/handlers"
    "github.com/holladworld/string-analyzer/database"
    "github.com/gin-gonic/gin"
)

func main() {
    // Initialize database
    err := database.Init()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    router := gin.Default()
    
    // Improved health check endpoint
    router.GET("/health", func(c *gin.Context) {
        // Simple response that always works
        c.JSON(200, gin.H{
            "status": "healthy",
            "service": "string-analyzer",
            "timestamp": "2024-01-21T10:00:00Z",
        })
    })
    
    // All required endpoints
    router.POST("/strings", handlers.PostStringHandler)
    router.GET("/strings/:string_value", handlers.GetStringHandler)
    router.GET("/strings", handlers.GetAllStringsHandler)
    router.GET("/strings/filter-by-natural-language", handlers.NaturalLanguageFilterHandler)
    router.DELETE("/strings/:string_value", handlers.DeleteStringHandler)
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    router.Run(":" + port)
}
