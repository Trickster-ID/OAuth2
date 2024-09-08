package model

type Api struct {
	StatusCode    int         `json:"status_code"`
	StatusMessage string      `json:"status_message"`
	Data          interface{} `json:"data,omitempty"`
	TotalCount    int         `json:"total_count,omitempty"`
	ErrorLog      *ErrorLog   `json:"error,omitempty"`
}
