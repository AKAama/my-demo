package models

// APIResponse 统一的API响应格式
type APIResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) *APIResponse {
	return &APIResponse{
		Status: 200,
		Data:   data,
		Msg:    "",
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(status int, msg string) *APIResponse {
	return &APIResponse{
		Status: status,
		Data:   nil,
		Msg:    msg,
	}
}
