package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockResponseHeaders struct {
	mock.Mock
}

func (r *MockResponseHeaders) writeHeaders(headers map[string][]string, res http.ResponseWriter) {
	r.Called(headers, res)
}

func Exchange_MockResponseHeaders(target *ResponseHeaders, mockVar *MockResponseHeaders) func() {
	var _headersMock ResponseHeaders = mockVar
	return exchange(target, &_headersMock)
}

type MockResponseWriter struct {
	mock.Mock
	header http.Header
}

func (i *MockResponseWriter) Header() http.Header {
	i.Called()
	return i.header
}

func (i *MockResponseWriter) Write(data []byte) (int, error) {
	i.Called(data)
	return 0, nil
}

func (i *MockResponseWriter) WriteHeader(statusCode int) {
	i.Called(statusCode)
}

func TestHeaders(t *testing.T) {
	responseWriter := MockResponseWriter{header: http.Header{}}
	headers := map[string][]string{"content-type": {"text/html", "utf-8"}}
	subject := ResponseHeadersWriter{}

	responseWriter.On("Header").Return(responseWriter.header).Twice()

	subject.writeHeaders(headers, &responseWriter)

	responseWriter.AssertExpectations(t)

	actualHeaders := map[string][]string(responseWriter.header)
	assert.Equal(t, map[string][]string{"Content-Type": {"text/html", "utf-8"}}, actualHeaders)
}
