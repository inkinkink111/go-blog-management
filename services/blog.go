package services

import (
	"context"
	"encoding/json"
	"fmt"
	"inkinkink111/go-blog-management/db"
	"inkinkink111/go-blog-management/models"
	"inkinkink111/go-blog-management/repositories"
	"inkinkink111/go-blog-management/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// @Summary Get all blogs
// @Description Get paginated list of blogs with optional tag filtering
// @Tags blogs
// @Accept json
// @Produce json
// @Param page query string false "Page number" default(1)
// @Param limit query string false "Items per page" default(10)
// @Param tags query string false "Comma-separated tags"
// @Success 200 {object} models.GetAllBlogRequest
// @Failure 500 {object} models.ResponseError
// @Router /api/v1/all_blogs [get]
func GetAllBlogs(c *fiber.Ctx) error {
	// Get query params
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")
	tags := c.Query("tags", "")
	// Convert page and limit to int
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}
	// Convert tags to slice
	var tagSlice []string
	if tags != "" {
		tagSlice = strings.Split(tags, ",")
		// Trim spaces from each tag
		for i, tag := range tagSlice {
			tagSlice[i] = strings.TrimSpace(tag)
		}
	}
	// Check for cache hit
	cacheKey := utils.GenerateCacheKey(page, limit, tagSlice)
	cachedResult := db.RedisClient.Get(context.Background(), cacheKey)
	if (cachedResult.Err() == nil) && (cachedResult.Val() != "") {
		var cachedResponse map[string]any
		err := json.Unmarshal([]byte(cachedResult.Val()), &cachedResponse)
		if err == nil {
			return c.Status(fiber.StatusOK).JSON(models.ResponseData{
				Message: "Get all blogs successfully.",
				Data:    cachedResponse,
			})
		}
	}
	// Cache miss - get from database
	blogRepo := repositories.NewBlogRepository()
	blogs, totalCount, err := blogRepo.GetAllBlogs(pageInt, limitInt, tagSlice)
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(models.ResponseError{
			Message: "Failed to get all blogs.",
			Error:   err.Error(),
		})
	}
	// Prep resp data
	respData := map[string]any{
		"blogs":       blogs,
		"page":        pageInt,
		"limit":       limitInt,
		"total_pages": (totalCount + int64(limitInt) - 1) / int64(limitInt),
		"total_item":  totalCount,
	}
	// Set cache
	cacheValue, _ := json.Marshal(respData)
	db.RedisClient.Set(context.Background(), cacheKey, cacheValue, 7*24*time.Hour)
	// Send response
	return c.Status(fiber.StatusOK).JSON(models.ResponseData{
		Message: "Get all blogs successfully.",
		Data:    respData,
	})
}

// @Summary Get blog by id
// @Tags blogs
// @Accept json
// @Produce json
// @Success 200 {object} models.GetBlogByIDResponse
// @Failure 500 {object} models.ResponseError
// @Router /api/v1/blogs/:blog_id [get]
func GetBlogByID(c *fiber.Ctx) error {
	// Get blog id
	blogID := c.Params("blog_id")
	// Validate
	if blogID == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(models.ResponseMsg{
			Message: "Missing blog id.",
		})
	}
	// Check cache
	cacheKey := fmt.Sprintf("blog:post:%s", blogID)
	cachedResult := db.RedisClient.Get(context.Background(), cacheKey)
	if (cachedResult.Err() == nil) && (cachedResult.Val() != "") {
		var blog models.Blog
		err := json.Unmarshal([]byte(cachedResult.Val()), &blog)
		if err != nil {
			return c.Status(fiber.ErrInternalServerError.Code).JSON(models.ResponseError{
				Message: "Failed to get blog.",
				Error:   err.Error(),
			})
		}
		println("Get blog from cache")
		return c.Status(fiber.StatusOK).JSON(models.ResponseData{
			Message: "Get blog successfully.",
			Data:    blog,
		})
	}
	// No cache hit, get blog from database
	blogRepo := repositories.NewBlogRepository()
	blog, err := blogRepo.GetBlogByID(blogID)
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(models.ResponseError{
			Message: "Failed to get blog.",
			Error:   err.Error(),
		})
	}
	if blog == nil {
		return c.Status(fiber.ErrNotFound.Code).JSON(models.ResponseMsg{
			Message: "Blog not found.",
		})
	}
	return c.Status(fiber.StatusOK).JSON(models.ResponseData{
		Message: "Get blog successfully.",
		Data:    blog,
	})
}

// @Summary Create a new blog post
// @Description Create a new blog post with title, content, and tags
// @Tags blogs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param blog body models.CreateBlogRequest true "Blog data"
// @Success 200 {object} models.CreateBlogSuccess
// @Failure 400 {object} models.CreateBlogError
// @Failure 500 {object} models.CreateBlogError
// @Router /api/v1/create_blog [post]
func CreateBlog(c *fiber.Ctx) error {
	// Get Author ID from jwt
	authorID := c.Locals("userId").(string)
	// Extract body
	body := &models.Blog{}
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(models.ResponseError{
			Message: "Invalid body.",
			Error:   err.Error(),
		})
	}
	// Validate
	if (body.Title == "") || (body.Content == "") || (authorID == "") || (len(body.Tags) == 0) {
		return c.Status(fiber.ErrBadRequest.Code).JSON(models.ResponseError{
			Message: "Invalid body.",
			Error:   "Missing required fields.",
		})
	}
	// Prep data
	body.Slug = utils.GenerateSlug(body.Title)
	body.AuthorID = authorID
	body.CreatedAt = time.Now()
	body.UpdatedAt = time.Now()
	body.BlogID = utils.GenerateID()
	// Create blog
	blogRepo := repositories.NewBlogRepository()
	err := blogRepo.InsertBlog(body)
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(models.ResponseError{
			Message: "Failed to create blog.",
			Error:   err.Error(),
		})
	}
	// Cache the newly created blog
	cacheKey := fmt.Sprintf("blog:post:%s", body.BlogID)
	cleanBody := models.Blog{
		BlogID:    body.BlogID,
		Title:     body.Title,
		Slug:      body.Slug,
		AuthorID:  body.AuthorID,
		Content:   body.Content,
		Tags:      body.Tags,
		CreatedAt: body.CreatedAt,
		UpdatedAt: body.UpdatedAt,
	}
	blogJSON, _ := json.Marshal(cleanBody)
	// 7 days cache
	db.RedisClient.Set(context.Background(), cacheKey, blogJSON, 7*24*time.Hour)
	// Invalidate list caches (since we added a new blog)
	pattern := "blog:list:*"
	keys, _ := db.RedisClient.Keys(context.Background(), pattern).Result()
	if len(keys) > 0 {
		db.RedisClient.Del(context.Background(), keys...)
	}

	return c.Status(fiber.StatusOK).JSON(models.ResponseData{
		Message: "Blog created successfully.",
		Data: map[string]string{
			"blog_id": body.BlogID,
		},
	})
}

// @Summary Update a blog post
// @Description Update a blog post with title, content, and tags
// @Tags blogs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param blog body models.CreateBlogRequest true "Blog data"
// @Success 200 {object} models.ResponseMsg
// @Failure 400 {object} models.ResponseError
// @Failure 404 {object} models.ResponseMsg
// @Failure 500 {object} models.ResponseError
// @Router /api/v1/update_blog/:blog_id [put]
func UpdateBlog(c *fiber.Ctx) error {
	// Get author id & blog id
	blogID := c.Params("blog_id")
	authorID := c.Locals("userId").(string)
	// Extract body
	body := &models.Blog{}
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(models.ResponseError{
			Message: "Invalid body.",
			Error:   err.Error(),
		})
	}
	// Validate
	if (body.Title == "") || (body.Content == "") || (authorID == "") || (len(body.Tags) == 0) {
		return c.Status(fiber.ErrBadRequest.Code).JSON(models.ResponseError{
			Message: "Invalid body.",
			Error:   "Missing required fields.",
		})
	}
	blogRepo := repositories.NewBlogRepository()
	// Check if blog exists
	blog, err := blogRepo.GetBlogByID(blogID)
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(models.ResponseError{
			Message: "Failed to find blog.",
			Error:   err.Error(),
		})
	}
	if blog == nil {
		return c.Status(fiber.ErrNotFound.Code).JSON(models.ResponseMsg{
			Message: "Blog not found.",
		})
	}
	// Check if blog is owned by user
	if blog.AuthorID != authorID {
		return c.Status(fiber.ErrForbidden.Code).JSON(models.ResponseMsg{
			Message: "You are not authorized to update this blog.",
		})
	}
	// Update blog
	updatedBlog := &models.Blog{
		BlogID:    blogID,
		Title:     body.Title,
		Content:   body.Content,
		Tags:      body.Tags,
		Slug:      utils.GenerateSlug(body.Title),
		CreatedAt: blog.CreatedAt,
		AuthorID:  blog.AuthorID,
		UpdatedAt: time.Now(),
	}
	err = blogRepo.UpdateBlog(updatedBlog)
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(models.ResponseError{
			Message: "Failed to update blog.",
			Error:   err.Error(),
		})
	}
	// Cache the updated blog
	cacheKey := fmt.Sprintf("blog:post:%s", blogID)
	updatedBlogJSON, _ := json.Marshal(updatedBlog)
	db.RedisClient.Set(context.Background(), cacheKey, updatedBlogJSON, 24*7*time.Hour)
	// Invalidate list caches (since we updated a blog)
	pattern := "blog:list:*"
	keys, _ := db.RedisClient.Keys(context.Background(), pattern).Result()
	if len(keys) > 0 {
		db.RedisClient.Del(context.Background(), keys...)
	}
	return c.Status(fiber.StatusOK).JSON(models.ResponseMsg{
		Message: "Blog updated successfully.",
	})
}

// @Summary Delete a blog post
// @Tags blogs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.ResponseMsg
// @Failure 400 {object} models.ResponseMsg
// @Failure 403 {object} models.ResponseMsg
// @Failure 404 {object} models.ResponseMsg
// @Failure 500 {object} models.ResponseError
// @Router /api/v1/delete_blog/:blog_id [delete]
func DeleteBlog(c *fiber.Ctx) error {
	// Get blog id and author id
	blogID := c.Params("blog_id")
	authorID := c.Locals("userId").(string)
	// Validate
	if blogID == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(models.ResponseMsg{
			Message: "Invalid body missing blog id.",
		})
	}
	blogRepo := repositories.NewBlogRepository()
	// Check if blog exists
	blog, err := blogRepo.GetBlogByID(blogID)
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(models.ResponseError{
			Message: "Failed to delete blog.",
			Error:   err.Error(),
		})
	}
	if blog == nil {
		return c.Status(fiber.ErrNotFound.Code).JSON(models.ResponseMsg{
			Message: "Blog not found.",
		})
	}
	// Check if blog is owned by user
	if blog.AuthorID != authorID {
		return c.Status(fiber.ErrForbidden.Code).JSON(models.ResponseMsg{
			Message: "You are not authorized to delete this blog.",
		})
	}
	// Delete blog
	err0 := blogRepo.DeleteBlog(blogID)
	if err0 != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(models.ResponseError{
			Message: "Failed to delete blog.",
			Error:   err.Error(),
		})
	}
	// Delete cache
	cacheKey := fmt.Sprintf("blog:post:%s", blogID)
	db.RedisClient.Del(context.Background(), cacheKey)
	// Invalidate list caches (since we deleted a blog)
	pattern := "blog:list:*"
	keys, _ := db.RedisClient.Keys(context.Background(), pattern).Result()
	if len(keys) > 0 {
		db.RedisClient.Del(context.Background(), keys...)
	}

	return c.Status(fiber.StatusOK).JSON(models.ResponseMsg{
		Message: "Blog deleted successfully.",
	})
}
