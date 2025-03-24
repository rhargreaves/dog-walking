package common

type ApiGatewayError struct {
	Message string `json:"message"`
}

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ApiErrorResponse struct {
	Error ApiError `json:"error"`
}
