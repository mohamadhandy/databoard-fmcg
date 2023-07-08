package helper

type Response struct {
	Data       any    `json:"data"`
	Error      error  `json:"error"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type ImageURL struct {
	ImageUrl   string `json:"image_url"`
	Error      error  `json:"error"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type FirebaseImage struct {
	Token string `json:"downloadTokens"`
}
