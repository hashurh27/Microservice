package forms

import (
	"Warranty-Microservice/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ClaimHandler(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		// Handle error from ConnectedDb()
		db, err := config.ConnectedDb()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// Get the data and handle the error
		data, err := config.ReadAllData(db, "claim")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// Set response header and write JSON data
		c.JSON(http.StatusOK, data)
	case http.MethodPost:
		// Parse the JSON request body into a map
		var requestData map[string]interface{}
		err := c.ShouldBindJSON(&requestData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request: Invalid JSON"})
			return
		}

		// Validate required fields in the JSON payload
		requiredFields := []string{"ClaimDate", "Description", "ClaimStatus"}
		for _, field := range requiredFields {
			if _, exists := requestData[field]; !exists {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Bad Request: Missing required field - %s", field)})
				return
			}
		}

		// Handle error from ConnectedDb()
		db, err := config.ConnectedDb()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// Insert the data into the database
		err = config.InsertData(db, "claim", requestData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// Respond with success message
		c.JSON(http.StatusCreated, gin.H{"message": "Data inserted successfully."})
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method Not Allowed"})
	}
}
