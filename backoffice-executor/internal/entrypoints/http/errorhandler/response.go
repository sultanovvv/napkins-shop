package errorhandler

// Response — общий контракт ошибки административного HTTP API.
type Response struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []Detail `json:"details,omitempty"`
}

// Detail описывает ошибку конкретного поля, например slug или parentId.
type Detail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Target  string `json:"target"`
}

func response(code, message string, details ...Detail) Response {
	return Response{Error: ErrorBody{Code: code, Message: message, Details: details}}
}
