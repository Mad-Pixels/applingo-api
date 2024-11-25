package v1

type APIResponse struct {
	Message string `json:"message,omitempty"`
}

var SuccessResponse = APIResponse{
	Message: "ok",
}
