package routes

import (
	"inkinkink111/go-blog-management/middleware"
	"inkinkink111/go-blog-management/services"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Blog API is running!"})
	})

	v1 := app.Group("/api/v1")
	v1.Post("/register", services.Register)
	v1.Post("/login", services.Login)
	v1.Get("/all_blogs", services.GetAllBlogs)
	v1.Get("/blog/:blog_id", services.GetBlogByID)

	auth := v1.Group("/")
	auth.Use(middleware.Authenticate)
	auth.Post("/create_blog", services.CreateBlog)
	auth.Delete("/delete_blog/:blog_id", services.DeleteBlog)
	auth.Put("/update_blog/:blog_id", services.UpdateBlog)
}
