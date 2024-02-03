package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type indicator struct {
    ElectricityStatus  bool  `json:"estatus"`
    LocationName  string  `json:"location_name"`
    LocationID string  `json:"location_id"`
    CurrentTime  string `json:"timestamp"`
}

var indicators = []indicator{
	{ElectricityStatus: true, LocationName: "Sompura Gate", LocationID: "562125", CurrentTime: "2024-02-03 14:20:23"},
}

func main() {
    router := gin.Default()
    router.GET("/status", getStatus)
    router.POST("/status", postStatus)

    router.Run("localhost:8080")
}

func getStatus(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, indicators)
}

func postStatus(c *gin.Context) {
    var newInd indicator

    if err := c.BindJSON(&newInd); err != nil {
        return
    }

    // Add the new ind to the slice.
    indicators = append(indicators, newInd)
    c.IndentedJSON(http.StatusCreated, newInd)
}


