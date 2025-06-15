package main

import (
	"github.com/gofiber/fiber/v2"
	"hiruna.com/project/config"
	"hiruna.com/project/routes"
	"hiruna.com/project/scripts"
)

func main() {
	app := fiber.New()
	config.ConnectDB()
	
	if config.AppSettings.ImportCSVOnStart {
		scripts.ImportSalesCSV(config.AppSettings.CSVFilePath)
	}
	config.Init()

	routes.RegisterRoutes(app)
	app.Listen(":3000")
}
