package response

type StreamResponse struct {
	code int
}

func (a StreamResponse) StatusCode() int {
	return a.code
}

func NewStreamResponse(code int) StreamResponse {
	return StreamResponse{code}
}
