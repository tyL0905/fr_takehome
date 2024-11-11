package main

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "net/http"
)

var pointsMap = make(map[string]int)

func main() {
    r := gin.Default()

    r.POST("/receipts/process", processReceipt)
    r.GET("/receipts/:id/points", getPoints)

    r.Run() 
}

func processReceipt(c *gin.Context) {
    var receipt Receipt

    if err := c.ShouldBindJSON(&receipt); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "The receipt is invalid"})
        return
    }

    points, err := calculatePoints(receipt)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "The receipt is invalid"})
        return
    }

    // Generate unique ID for the receipt
    receipt.ID = uuid.New().String()
    
    // save points to map
    pointsMap[receipt.ID] = points

    c.JSON(http.StatusOK, gin.H{"id": receipt.ID})
}


func getPoints(c *gin.Context) {
    id := c.Param("id")
    points, exists := pointsMap[id]

    if !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "No receipt found for that id"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"points": points})
}
