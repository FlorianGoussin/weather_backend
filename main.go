package main

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "github.com/weather-backend/models"
)

func main() {
  models.mongoConnect()


  r := gin.Default()

  r.GET("/", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"data": "Root route"})    
  })

  r.GET("/cities", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"data": "Cities autocomplete"})    
  })

  r.Run()
}