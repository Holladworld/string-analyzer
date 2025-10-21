package main

import (
    "log"
    "os"
    "string-analyzer/handlers"
    "string-analyzer/database"
    "github.com/gin-gonic/gin"
)

func main() {
    // Initialize database
    err := database.Init()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    router := gin.Default()
    
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
