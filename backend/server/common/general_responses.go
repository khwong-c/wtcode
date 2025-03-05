package common

type GeneralResponse struct {
	Success bool `json:"success"`
}

var (
	SuccessResponse = GeneralResponse{Success: true}
)
