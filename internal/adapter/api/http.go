package api

import (
	ports "jamo/backend/internal/port"
)

type HTTPHandler struct {
	Service ports.Service
}

func NewHTTPHandler(service ports.Service) *HTTPHandler {
	return &HTTPHandler{
		Service: service,
	}
}
