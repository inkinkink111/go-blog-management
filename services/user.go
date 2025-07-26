package services

import (
	"time"

	"inkinkink111/go-blog-management/models"
	"inkinkink111/go-blog-management/repositories"
	"inkinkink111/go-blog-management/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// @Summary Register user
// @Tags auth
// @Accept json
// @Produce json
// @Param register body object{email=string,password=string,name=string} true "User data"
// @Success 200 {object} models.ResponseMsg
// @Failure 400 {object} models.ResponseError
// @Failure 409 {object} models.ResponseMsg
// @Failure 500 {object} models.ResponseError
// @Router /api/v1/register [post]
func Register(c *fiber.Ctx) error {
	// Extract body
	body := &models.User{}
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(models.ResponseError{
			Message: "Invalid body.",
			Error:   err.Error(),
		})
	}
	// Validate
	if (body.Email == "") || (body.Password == "") || (body.Name == "") {
		return c.Status(fiber.ErrBadRequest.Code).JSON(models.ResponseError{
			Message: "Invalid body.",
			Error:   "Missing required fields.",
		})
	}
	// Check if user already exists
	userRepo := repositories.NewUserRepository()
	existingUser, err := userRepo.GetUserByEmail(body.Email)
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(models.ResponseError{
			Message: "Internal server error.",
			Error:   err.Error(),
		})
	}
	if existingUser != nil {
		return c.Status(fiber.ErrConflict.Code).JSON(models.ResponseMsg{
			Message: "User already exists",
		})
	}
	// Hash password
	hashedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(models.ResponseError{
			Message: "Internal server error.",
			Error:   err.Error(),
		})
	}
	// Store in Mongo
	body.Password = hashedPassword
	body.CreatedAt = time.Now()
	body.UserId = uuid.NewString()

	if err := userRepo.InsertUser(body); err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(models.ResponseError{
			Message: "Internal server error.",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.ResponseMsg{
		Message: "Create user successfully.",
	})
}

// @Summary Login
// @Tags auth
// @Accept json
// @Produce json
// @Param login body object{email=string,password=string} true "User credentials"
// @Success 200 {object} object{message=string,data=object{token=string}}
// @Failure 400 {object} models.ResponseError
// @Failure 401 {object} models.ResponseMsg
// @Failure 404 {object} models.ResponseMsg
// @Failure 500 {object} models.ResponseError
// @Router /api/v1/login [post]
func Login(c *fiber.Ctx) error {
	body := &models.User{}
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(models.ResponseError{
			Message: "Invalid body.",
			Error:   err.Error(),
		})
	}
	// Validate
	if (body.Email == "") || (body.Password == "") {
		return c.Status(fiber.ErrBadRequest.Code).JSON(models.ResponseError{
			Message: "Invalid body.",
			Error:   "Missing required fields.",
		})
	}
	// Get user
	userRepo := repositories.NewUserRepository()
	user, err := userRepo.GetUserByEmail(body.Email)
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(models.ResponseError{
			Message: "Failed to get user.",
			Error:   err.Error(),
		})
	}
	if user == nil {
		return c.Status(fiber.ErrNotFound.Code).JSON(models.ResponseMsg{
			Message: "Invalid email or password.",
		})
	}
	// Compare Password
	if !utils.ComparePassword(user.Password, body.Password) {
		return c.Status(fiber.ErrUnauthorized.Code).JSON(models.ResponseMsg{
			Message: "Invalid email or password.",
		})
	}
	// Generate token and return to client
	token, err := utils.GenerateToken(user.Email, user.UserId)
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(models.ResponseError{
			Message: "Failed to generate token.",
			Error:   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(models.ResponseData{
		Message: "Login successfully.",
		Data: map[string]string{
			"token": token,
		},
	})
}
