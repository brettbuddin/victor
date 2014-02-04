package httpserver

import (
	"net/http"
)

type ResponseWriteTracker struct {
	http.ResponseWriter
	code int
	size int64
}

func (w *ResponseWriteTracker) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *ResponseWriteTracker) Write(p []byte) (int, error) {
	if w.code == 0 {
		w.code = 200
	}

	w.size += int64(len(p))
	return w.ResponseWriter.Write(p)
}
