package server

import (
	"net/http"
	"testing"
)

func Test_ResponseBodyString(t *testing.T) {
	responseWriter := MockResponseWriter{header: http.Header{}}
	data := []byte{0, 1, 2, 3}
	headers := map[string][]string{"h1": {"v1", "v2"}}
	subject := ResponseBodyString{response: data, Response: Response{statusCode: 200, headers: headers}}

	var mockHeaders *MockResponseHeaders = &MockResponseHeaders{}
	defer Exchange_MockResponseHeaders(&responseHeadersWriter, mockHeaders)()

	responseWriter.On("Write", data).Return(0, nil)
	mockHeaders.On("writeHeaders", headers, &responseWriter).Once()
	responseWriter.On("WriteHeader", 200).Once()

	subject.WriteTo(&responseWriter)

	responseWriter.AssertExpectations(t)
	mockHeaders.AssertExpectations(t)
}
