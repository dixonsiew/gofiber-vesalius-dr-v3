package auth

import (
    "vesaliusdr/dto"
    doctorAppUserService "vesaliusdr/service/doctor_app_user"
    "vesaliusdr/utils"

    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v3"
)

// Login
//
// @Tags Auth
// @Produce json
// @Param request body dto.LoginDto true "Login Request"
// @Success 200
// @Router /login [post]
func Login(c fiber.Ctx) error {
    data := new(dto.LoginDto)
    mx := fiber.Map{
        "statusCode": fiber.StatusUnauthorized,
        "message":    "Invalid Credentials",
    }
    if err := c.Bind().Body(data); err != nil {
        if validationErrors, ok := err.(validator.ValidationErrors); ok {
            errs := utils.GetValidationErrors(validationErrors)
            if errs != nil {
                return c.Status(fiber.StatusUnauthorized).JSON(mx)
            }
        }

        return c.Status(fiber.StatusUnauthorized).JSON(mx)
    }

    if data.FromBiometric == 1 {
        user, err := doctorAppUserService.FindByUsername(data.Username)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(mx)
        }

        valid := false
        if user != nil {
            a := *user
            if a.Machine_id != "" {
                valid = doctorAppUserService.ValidateCredentials2(a, data.Password)
                if !valid {
                    return c.Status(fiber.StatusUnauthorized).JSON(mx)
                }

                token, err := utils.GenerateToken(a)
                refreshToken, errx := utils.GenerateRefreshToken(a)
                if err != nil {
                    return c.Status(fiber.StatusUnauthorized).JSON(mx)
                }

                if errx != nil {
                    return c.Status(fiber.StatusUnauthorized).JSON(mx)
                }

                c.Set(fiber.HeaderAuthorization, token)
                return c.JSON(fiber.Map{
                    "type":             "bearer",
                    "token":            token,
                    "refresh_token":    refreshToken,
                    "isFirstTimeLogin": a.First_time_login,
                    "role":             a.Role,
                    "mcr":              a.Mcr,
                    "branch":           a.Branch,
                })
            } else {
                return c.Status(fiber.StatusUnauthorized).JSON(mx)
            }
        } else {
            return c.Status(fiber.StatusUnauthorized).JSON(mx)
        }
    } else {
        user, err := doctorAppUserService.FindByUsername(data.Username)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(mx)
        }

        valid := false
        if user != nil {
            valid = doctorAppUserService.ValidateCredentials(*user, data.Password)
        }

        if !valid {
            return c.Status(fiber.StatusUnauthorized).JSON(mx)
        }

        a := *user
        token, err := utils.GenerateToken(a)
        refreshToken, errx := utils.GenerateRefreshToken(a)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(mx)
        }

        if errx != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(mx)
        }

        c.Set(fiber.HeaderAuthorization, token)
        return c.JSON(fiber.Map{
            "type":             "bearer",
            "token":            token,
            "refresh_token":    refreshToken,
            "isFirstTimeLogin": a.First_time_login,
            "role":             a.Role,
            "mcr":              a.Mcr,
            "branch":           a.Branch,
        })
    }
}
