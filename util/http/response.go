package http

import "net/http"

type Response struct {
	body []byte
	raw  *http.Response
}

func newResponse(body []byte, raw *http.Response) *Response {
	return &Response{body: body, raw: raw}
}

func (r *Response) Body() []byte {
	return r.body
}

func (r *Response) Raw() *http.Response {
	return r.raw
}
