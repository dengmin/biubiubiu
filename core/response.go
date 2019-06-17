package core

import (
	"gopkg.in/djherbis/stream.v1"
	"net/http"
)

func newResponseStreamer(w http.ResponseWriter) *ResponseStreamer {
	strm, err := stream.NewStream("responseBuffer", stream.NewMemFS())
	if err != nil {
		panic(err)
	}
	return &ResponseStreamer{
		ResponseWriter: w,
		Stream:         strm,
		C:              make(chan struct{}),
	}
}

type ResponseStreamer struct {
	StatusCode int
	http.ResponseWriter
	*stream.Stream
	C chan struct{}
}

func (rw *ResponseStreamer) WaitHeaders() {
	for range rw.C {
	}
}

func (rw *ResponseStreamer) WriteHeader(status int) {
	defer close(rw.C)
	rw.StatusCode = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *ResponseStreamer) Write(b []byte) (int, error) {
	rw.Stream.Write(b)
	return rw.ResponseWriter.Write(b)
}
func (rw *ResponseStreamer) Close() error {
	return rw.Stream.Close()
}
