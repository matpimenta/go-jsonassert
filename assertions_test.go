package jsonassert

import (
    "fmt"
    "testing"
    "github.com/stretchr/testify/assert"
)

type MockTesting struct {
    Failed bool
    Message string
}

func (t *MockTesting) Errorf(format string, args ...interface{}) {
    t.Failed = true
    t.Message = fmt.Sprintf(format, args)
}

func TestAssertJSONEqualsWhenJSONsAreEqual(t *testing.T) {
    tt := &MockTesting{}
    AssertJSONEquals(tt, `{"url": "URL 1"}`, `{"url": "URL 1"}`, false)
    assert.False(t, tt.Failed, tt.Message)
}

func TestAssertJSONEqualsWhenJSONsAreNotEqual(t *testing.T) {
    tt := &MockTesting{}
    AssertJSONEquals(tt, `{"url": "URL 1"}`, `{"url": "URL 2"}`, false)
    assert.True(t, tt.Failed, tt.Message)
    assert.Equal(t, "[url:\nExpected: \"URL 1\"\ngot: \"URL 2\"\n]", tt.Message)
}

func TestAssertJSONNotEqualsWhenJSONsAreEqual(t *testing.T) {
    tt := &MockTesting{}
    AssertJSONEquals(tt, `{"url": "URL 1"}`, `{"url": "URL 1"}`, false)
    assert.False(t, tt.Failed, tt.Message)
}

func TestAssertJSONNotEqualsWhenJSONsAreNotEqual(t *testing.T) {
    tt := &MockTesting{}
    AssertJSONEquals(tt, `{"url": "URL 1"}`, `{"url": "URL 2"}`, false)
    assert.True(t, tt.Failed, tt.Message)
}
