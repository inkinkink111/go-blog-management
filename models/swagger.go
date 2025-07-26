package models

type CreateBlogRequest struct {
	Title   string   `json:"title" example:"My Blog Title" validate:"required"`
	Content string   `json:"content" example:"Blog content" validate:"required"`
	Tags    []string `json:"tags" example:"golang,redis" validate:"required"`
}

type CreateBlogSuccess struct {
	Message string            `json:"message" example:"Blog created successfully."`
	Data    map[string]string `json:"data" example:"blog_id: 1234567890"`
}

type CreateBlogError struct {
	Message string `json:"message" example:"Failed to create blog."`
	Error   string `json:"error" example:"Failed to create blog."`
}

type GetAllBlogRequest struct {
	Message string            `json:"message" example:"Get all blogs successfully."`
	Data    map[string]string `json:"data" example:"blogs:blog_data,page:1,limit:10,total_pages:1,total_item:1"`
}

type GetBlogByIDResponse struct {
	Message string `json:"message" example:"Get blog by id successfully."`
	Data    Blog   `json:"data" example:{
		"id": "1234567890",
		"blog_id": "1234567890",
		"title": "My Blog Title",
		"content": "Blog content",
		"slug": "my-blog-title",
		"created_at": "2021-01-01T00:00:00Z",
		"tags": ["golang", "redis"],
		"author_id": "1234567890",
		"updated_at": "2021-01-01T00:00:00Z"
	}`
}
