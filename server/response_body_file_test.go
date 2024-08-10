package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func exchange[T any](first *T, second *T) func() {
	*first, *second = *second, *first
	return func() {
		*first, *second = *second, *first
	}
}

func Test_ResponseBodyFile(t *testing.T) {
	// headers
	var mockHeaders *MockResponseHeaders = &MockResponseHeaders{}
	defer Exchange_MockResponseHeaders(&responseHeadersWriter, mockHeaders)()

	// file utils
	mockFileUtils := &MockFileUtils{data: []byte{1, 2, 3, 4, 5}}
	defer Exchange_MockFileUtils(&fileUtils, mockFileUtils)()

	responseWriter := MockResponseWriter{header: http.Header{}}

	headers := map[string][]string{"h1": {"v1", "v2"}}
	subject := ResponseBodyFile{file: "test-name", Response: Response{statusCode: 200, headers: headers}}

	mockHeaders.On("writeHeaders", headers, &responseWriter).Once()
	mockFileUtils.On("read", "test-name", mock.Anything).Once()
	responseWriter.On("Write", mockFileUtils.data).Once()
	responseWriter.On("WriteHeader", 200).Once()

	subject.WriteTo(&responseWriter)

	responseWriter.AssertExpectations(t)
	mockFileUtils.AssertExpectations(t)
	mockHeaders.AssertExpectations(t)
}

func Test_ResponseBodyFileCachable(t *testing.T) {
	fileData := []byte{1, 2, 3, 4, 5}
	// headers
	var mockHeaders *MockResponseHeaders = &MockResponseHeaders{}
	defer Exchange_MockResponseHeaders(&responseHeadersWriter, mockHeaders)()

	// file utils
	mockFileUtils := &MockFileUtils{data: fileData}
	defer Exchange_MockFileUtils(&fileUtils, mockFileUtils)()

	responseWriter := MockResponseWriter{header: http.Header{}}

	headers := map[string][]string{"h1": {"v1", "v2"}}
	subject := ResponseBodyFileCachable{ResponseBodyFile: ResponseBodyFile{file: "test-name", Response: Response{statusCode: 200, headers: headers}}}

	mockHeaders.On("writeHeaders", headers, &responseWriter).Twice()
	mockFileUtils.On("read", "test-name", mock.Anything).Once()
	responseWriter.On("Write", fileData).Twice()
	responseWriter.On("WriteHeader", 200).Twice()

	assert.False(t, subject.isCached)
	assert.Equal(t, 0, len(subject.cache))

	subject.WriteTo(&responseWriter)

	assert.True(t, subject.isCached)
	assert.Equal(t, 5, len(subject.cache))

	// second call reads only cached data, no interaction with storage
	subject.WriteTo(&responseWriter)

	responseWriter.AssertExpectations(t)
	mockFileUtils.AssertExpectations(t)
	mockHeaders.AssertExpectations(t)
}
