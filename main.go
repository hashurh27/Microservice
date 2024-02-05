package main

import (
	"Warranty-Microservice/api/forms"
	"github.com/gin-gonic/gin"
	"log"
)

// ... rest of the code ...
func setupRoutes(router *gin.Engine) {
	apiGroup := router.Group("/api")

	customerGroup := apiGroup.Group("/customer")
	{
		customerGroup.GET("", gin.HandlerFunc(forms.CustomerHandler))
		customerGroup.POST("", gin.HandlerFunc(forms.CustomerHandler))
	}

	usersGroup := apiGroup.Group("/users")
	{
		usersGroup.GET("", forms.UserHandler)  // Remove gin.HandlerFunc conversion
		usersGroup.POST("", forms.UserHandler) // Remove gin.HandlerFunc conversion
	}

	productGroup := apiGroup.Group("/product")
	{
		productGroup.GET("", forms.ProductHandler)  // Remove gin.HandlerFunc conversion
		productGroup.POST("", forms.ProductHandler) // Remove gin.HandlerFunc conversion
	}

	warrantyGroup := apiGroup.Group("/warranty")
	{
		warrantyGroup.GET("", forms.WarrantyHandler)  // Ensure WarrantyHandler is defined in forms package
		warrantyGroup.POST("", forms.WarrantyHandler) // Ensure WarrantyHandler is defined in forms package
	}

	claimGroup := apiGroup.Group("/claim")
	{
		claimGroup.GET("", forms.ClaimHandler)  // Remove gin.HandlerFunc conversion
		claimGroup.POST("", forms.ClaimHandler) // Remove gin.HandlerFunc conversion
	}
}

func main() {
	// Create a new Gin router
	router := gin.Default()

	// Set up routes
	setupRoutes(router)

	// Start the server
	err := router.Run(":8080")
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
