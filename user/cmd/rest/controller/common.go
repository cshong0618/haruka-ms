package controller

type Response struct {
	Data  interface{} `json:"data"`
	Error *Error      `json:"error"`
}

type Error struct {
	Reason string `json:"reason"`
}

func NewErrorResponse(err error) Response {
	return Response{
		Data:  nil,
		Error: &Error{
			Reason: err.Error(),
		},
	}
}

func NewResponse(data interface{}) Response {
	return Response {
		Data: data,
		Error: nil,
	}
}