package user

import (
	"vesaliusdr/controller/user"
	"vesaliusdr/middleware"

	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    api := router.Group("/user")
    api.Use(middleware.JWTProtected)
    api.Get("/userId/:userId", user.GetUserById)
    api.Post("/first-time-change-password", user.FirstTimeChangePassword)
    api.Post("/change-password", user.ChangePassword)
    api.Post("/add-machine-id", user.AddMachineId)
    api.Get("/email/:email", user.GetUserByEmail)
    api.Get("", user.GetUser)
}