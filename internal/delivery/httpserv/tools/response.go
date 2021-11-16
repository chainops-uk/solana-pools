package tools

type ResponseData struct {
	Data interface{} `json:"data"`
}

type ResponseError struct {
	Error string `json:"error"`
}
