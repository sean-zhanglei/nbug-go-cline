package main

import (
	"mongo-echo-go/handler"
	"mongo-echo-go/utils"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize database
	utils.InitDB()

	e := echo.New()

	// Routes
	e.GET("/test", handler.GetTest)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
