package user

import (
    "strconv"
    "vesaliusdr/dto"
    "vesaliusdr/middleware"
    _ "vesaliusdr/model"
    doctorAppUserService "vesaliusdr/service/doctor_app_user"
    "vesaliusdr/utils"

    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v3"
)

// GetUserById
//
// @Tags User
// @Produce json
// @Param        userId   path      string  true  "UserId"
// @Security BearerAuth
// @Success 200 {object} model.DoctorAppUser
// @Router /user/userId/{userId} [get]
func GetUserById(c fiber.Ctx) error {
    userIds := c.Params("userId")
    userId, err := strconv.Atoi(userIds)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "statusCode": fiber.StatusNotFound,
            "message":    "Doctor Not Found",
        })
    }

    user, err := doctorAppUserService.FindByUserId(userId)
    if err != nil || user == nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "statusCode": fiber.StatusNotFound,
            "message":    "Doctor Not Found",
        })
    }

    return c.JSON(user)
}

// GetUserById
//
// @Tags User
// @Produce json
// @Param request body dto.PostChangePasswordDto true "ChangePassword Request"
// @Security BearerAuth
// @Success 200
// @Router /user/first-time-change-password [post]
func FirstTimeChangePassword(c fiber.Ctx) error {
    _, a, err := middleware.ValidateToken(c)
    if err != nil {
        return middleware.Unauthorized(c)
    }

    v := *a
    data := new(dto.PostChangePasswordDto)
    if err := c.Bind().Body(data); err != nil {
        if validationErrors, ok := err.(validator.ValidationErrors); ok {
            errs := utils.GetValidationErrors(validationErrors)
            if errs != nil {
                return middleware.Unauthorized(c)
            }
        }

        return middleware.Unauthorized(c)
    }

    valid := doctorAppUserService.ValidateCredentials(v, data.OldPassword)
    if !valid {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "statusCode": fiber.StatusBadRequest,
            "message":    "Current password is invalid",
        })
    }

    valid1 := doctorAppUserService.ValidateCredentials(v, data.NewPassword)
    if !valid1 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "statusCode": fiber.StatusBadRequest,
            "message":    "New password is not allow to be the same",
        })
    }

    v.Password = data.NewPassword
    v.First_time_login = false
    err = doctorAppUserService.FirstTimeChangePassword(v)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "statusCode": fiber.StatusInternalServerError,
            "message":    err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "successMessage": "Password has been updated",
    })
}

// ChangePassword
//
// @Tags User
// @Produce json
// @Param request body dto.PostChangePasswordDto true "ChangePassword Request"
// @Security BearerAuth
// @Success 200
// @Router /user/change-password [post]
func ChangePassword(c fiber.Ctx) error {
    _, a, err := middleware.ValidateToken(c)
    if err != nil {
        return middleware.Unauthorized(c)
    }

    v := *a
    data := new(dto.PostChangePasswordDto)
    if err := c.Bind().Body(data); err != nil {
        if validationErrors, ok := err.(validator.ValidationErrors); ok {
            errs := utils.GetValidationErrors(validationErrors)
            if errs != nil {
                return middleware.Unauthorized(c)
            }
        }

        return middleware.Unauthorized(c)
    }

    valid := doctorAppUserService.ValidateCredentials(v, data.OldPassword)
    if !valid {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "statusCode": fiber.StatusBadRequest,
            "message":    "Current password is invalid",
        })
    }

    valid1 := doctorAppUserService.ValidateCredentials(v, data.NewPassword)
    if !valid1 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "statusCode": fiber.StatusBadRequest,
            "message":    "New password is not allow to be the same",
        })
    }

    v.Password = data.NewPassword
    v.First_time_login = false
    err = doctorAppUserService.FirstTimeChangePassword(v)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "statusCode": fiber.StatusInternalServerError,
            "message":    err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "successMessage": "Password has been updated",
    })
}

// AddMachineId
//
// @Tags User
// @Produce json
// @Param request body dto.PostMachineInfo true "AddMachineId Request"
// @Security BearerAuth
// @Success 200
// @Router /user/add-machine-id [post]
func AddMachineId(c fiber.Ctx) error {
    _, a, err := middleware.ValidateToken(c)
    if err != nil {
        return middleware.Unauthorized(c)
    }

    v := *a
    data := new(dto.PostMachineInfo)
    if err := c.Bind().Body(data); err != nil {
        if validationErrors, ok := err.(validator.ValidationErrors); ok {
            errs := utils.GetValidationErrors(validationErrors)
            if errs != nil {
                return middleware.Unauthorized(c)
            }
        }

        return middleware.Unauthorized(c)
    }

    v.Machine_id = data.MachineId
    err = doctorAppUserService.UpdateMachineId(v)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "statusCode": fiber.StatusInternalServerError,
            "message":    err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "successMessage": "Machine Info has been added",
    })
}

// GetUserByEmail
//
// @Tags User
// @Produce json
// @Param        email   path      string  true  "Email"
// @Security BearerAuth
// @Success 200 {object} model.DoctorAppUser
// @Router /user/email/{email} [get]
func GetUserByEmail(c fiber.Ctx) error {
    email := c.Params("email")
    user, err := doctorAppUserService.FindByEmail(email)
    if err != nil || user == nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "statusCode": fiber.StatusNotFound,
            "message":    "Doctor Not Found",
        })
    }

    return c.JSON(user)
}

// GetUser
//
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.DoctorAppUser
// @Router /user [get]
func GetUser(c fiber.Ctx) error {
    _, user, ret := middleware.ValidateToken(c)
    if ret != nil {
        return middleware.Unauthorized(c)
    }

    return c.JSON(user)
}