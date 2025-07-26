package models

type ResponseMsg struct {
	Message string `json:"message"`
}

type ResponseData struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ResponseError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
