package auth

import (
    "vesaliusdr/controller/auth"
    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    router.Post("/login", auth.Login)
}