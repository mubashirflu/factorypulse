package machines

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateMachineHandler(c *gin.Context) {
	var input CreateMachineInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := CreateMachine(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "machine created", "machine_id": id})
}

func GetAllMachinesHandler(c *gin.Context) {
	list, err := GetAllMachines()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func GetMachineHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid machine id"})
		return
	}

	machine, err := GetMachineByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "machine not found"})
		return
	}

	c.JSON(http.StatusOK, machine)
}
