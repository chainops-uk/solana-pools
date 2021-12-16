package tools

type ResponseData struct {
	Data interface{} `json:"data"`
}

type ResponseArrayData struct {
	Data     interface{} `json:"data"`
	MetaData *MetaData   `json:"meta_data"`
}

type MetaData struct {
	Offset      uint64 `json:"offset"`
	Limit       uint64 `json:"limit"`
	TotalAmount uint64 `json:"total_amount"`
}

type ResponseError struct {
	Error string `json:"error"`
}
