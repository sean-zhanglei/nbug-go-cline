package main

import (
	"mongo-echo-go/handler"
	"mongo-echo-go/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	// Initialize database
	utils.InitDB()

	e := echo.New()

	// Routes
	e.GET("/test", handler.GetTest)
	e.GET("/testSum", handler.GetTestSum)

	// Start server
	e.Logger.SetLevel(log.DEBUG)
	e.Logger.Fatal(e.Start(":8080"))
}
