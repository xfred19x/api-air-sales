package main

import (
	"api-air-sales/configs"
	"api-air-sales/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	// Run Database
	configs.ConnectDB()

	// Routes
	routes.SetupRoutes(router) //add this

	router.Run("localhost:6000")
}
