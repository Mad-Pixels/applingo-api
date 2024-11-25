package openapi

// Response стандартная структура ответа
type Response struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

var (
	// DefaultSuccessResponse стандартный ответ для успешных операций
	DefaultSuccessResponse = Response{
		Message: "ok",
	}
)
