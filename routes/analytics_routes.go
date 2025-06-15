package routes

import (
	"github.com/gofiber/fiber/v2"
	"hiruna.com/project/controllers"
)

func RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/revenue/countries", controllers.GetChunkedViewData)
	api.Get("/products/top", controllers.GetTopProducts)
	api.Get("/sales/monthly", controllers.GetMonthlySales)
	api.Get("/revenue/regions", controllers.GetTopRegions)
}
