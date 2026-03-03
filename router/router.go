package router

import (
    "vesaliusdr/router/auth"
    "vesaliusdr/router/user"
    "vesaliusdr/router/vesalius"
    "github.com/gofiber/fiber/v3"
    // "github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
    api := app.Group("/mobile_central_dr-1.0.0")
    auth.SetupRoutes(api)
    user.SetupRoutes(api)
    vesalius.SetupRoutes(api)
}