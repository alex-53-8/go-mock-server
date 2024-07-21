package server

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFileUtils struct {
	mock.Mock
	data []byte
}

func (f *MockFileUtils) read(fileName string, consumer func(*[]byte)) error {
	f.Called(fileName, consumer)
	consumer(&f.data)
	return nil
}

func Exchange_MockFileUtils(target *FileUtils, mockVar *MockFileUtils) func() {
	var _mockFileUtils FileUtils = mockVar
	return exchange(&fileUtils, &_mockFileUtils)
}

func Test_FileUtilsService(t *testing.T) {
	subject := FileUtilsService{}

	// no file, then error is returned
	err := subject.read("/tmp/i-do-not-exist", func(*[]byte) {})
	assert.NotNil(t, err)

	f, err := os.CreateTemp("/tmp", "prefix")
	if err != nil {
		assert.Fail(t, err.Error())
	}

	data := []byte{1, 2, 3, 4, 5}

	size, err := f.Write(data)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, 5, size)

	actualData := []byte{}
	subject.read(f.Name(), func(d *[]byte) {
		actualData = append(actualData, *d...)
	})

	assert.Equal(t, data, actualData)
}

func Test_getCacheableState(t *testing.T) {
	file, err := os.CreateTemp("/tmp", "mock-server-test")
	defer os.Remove(file.Name())

	if err != nil {
		assert.Fail(t, err.Error())
	}

	_, err = file.Write([]byte{1, 2, 3, 4, 5})

	if err != nil {
		assert.Fail(t, err.Error())
	}

	// max items size stored in cache is greater than size of our file, can be cached
	assert.Equal(t, canBeCached, getCacheableState(file.Name(), 6))
	// max items size stored in cache equals to size of our file, can be cached
	assert.Equal(t, canBeCached, getCacheableState(file.Name(), 5))
	// max items size stored in cache is less than size of our file, cannot be cached
	assert.Equal(t, cannotBeCached, getCacheableState(file.Name(), 4))
}
