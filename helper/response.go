package helper

type Response struct {
	Data       any    `json:"data"`
	Error      error  `json:"error"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}
