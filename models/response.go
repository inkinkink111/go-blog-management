package models

type ResponseMsg struct {
	Message string `json:"message" example:"Response message"`
}

type ResponseData struct {
	Message string `json:"message" example:"Success"`
	Data    any    `json:"data"`
}

type ResponseError struct {
	Message string `json:"message" example:"Error message"`
	Error   string `json:"error" example:"Detailed error"`
}
