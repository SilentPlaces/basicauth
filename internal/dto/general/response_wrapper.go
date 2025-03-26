package general_reponse_dto

type (
	Response struct {
		Status  string      `json:"status"` // "success" or "error"
		Data    interface{} `json:"data,omitempty"`
		Message string      `json:"message,omitempty"`
		Code    int         `json:"code,omitempty"`
	}
)
