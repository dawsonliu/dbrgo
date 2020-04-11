package core

type StandardResult struct {
	StatusCode   int         `json:"statusCode"`
	ErrorMessage string      `json:"errorMessage"`
	Data         interface{} `json:"data"`
}

func SetResult(data interface{}) *StandardResult {
	return &StandardResult{
		StatusCode:   200,
		ErrorMessage: "",
		Data:         data,
	}
}

func SetError(errorMessage string, statusCode int) *StandardResult {
	return &StandardResult{
		StatusCode:   statusCode,
		ErrorMessage: errorMessage,
		Data:         nil,
	}
}

func SetDefaultError(errorMessage string) *StandardResult {
	return SetError(errorMessage, 500)
}
