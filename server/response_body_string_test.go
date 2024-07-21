package server

import (
	"net/http"
	"testing"
)

func Test_ResponseBodyString(t *testing.T) {
	responseWriter := MockResponseWriter{header: http.Header{}}
	data := []byte{0, 1, 2, 3}
	headers := map[string][]string{"h1": {"v1", "v2"}}
	subject := ResponseBodyString{response: data, headers: headers}

	var mockHeaders *MockResponseHeaders = &MockResponseHeaders{}
	defer Exchange_MockResponseHeaders(&responseHeadersWriter, mockHeaders)()

	responseWriter.On("Write", data).Return(0, nil)
	mockHeaders.On("writeHeaders", headers, &responseWriter).Return(true).Once()

	subject.WriteTo(&responseWriter)

	responseWriter.AssertExpectations(t)
	mockHeaders.AssertExpectations(t)
}
