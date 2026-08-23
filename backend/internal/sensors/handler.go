package sensors

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateReadingHandler(c *gin.Context) {
	var input CreateReadingInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := CreateReading(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "reading recorded"})
}

func GetLatestReadingHandler(c *gin.Context) {
	machineID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid machine id"})
		return
	}

	reading, err := GetLatestReading(machineID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no readings found"})
		return
	}

	c.JSON(http.StatusOK, reading)
}

func GetReadingHistoryHandler(c *gin.Context) {
	machineID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid machine id"})
		return
	}

	// Default 50 readings, agar query param diya ho toh wo use karo
	limit := 50
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	readings, err := GetReadingHistory(machineID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, readings)
}
