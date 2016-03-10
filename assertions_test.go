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

func TestAssertJSONLenientlyEqualsWhenJSONsAreEqual(t *testing.T) {
    tt := &MockTesting{}
    AssertJSONEquals(tt, `{"url": "URL 1"}`, `{"url": "URL 1", "name": "John"}`, false)
    assert.False(t, tt.Failed, tt.Message)
}

func TestAssertJSONLenientlyEqualsWhenJSONsAreNotEqual(t *testing.T) {
    tt := &MockTesting{}
    AssertJSONEquals(tt, `{"url": "URL 1"}`, `{"url": "URL 2", "name": "John"}`, false)
    assert.True(t, tt.Failed, tt.Message)
    assert.Equal(t, "[url:\nExpected: \"URL 1\"\ngot: \"URL 2\"\n]", tt.Message)
}

func TestAssertJSONLenientlyNotEqualsWhenJSONsAreEqual(t *testing.T) {
    tt := &MockTesting{}
    AssertJSONNotEquals(tt, `{"url": "URL 1"}`, `{"url": "URL 2"}`, false)
    assert.False(t, tt.Failed, tt.Message)
}

func TestAssertJSONLenientlyNotEqualsWhenJSONsAreNotEqual(t *testing.T) {
    tt := &MockTesting{}
    AssertJSONNotEquals(tt, `{"url": "URL 1"}`, `{"url": "URL 1"}`, false)
    assert.True(t, tt.Failed, tt.Message)
}

func TestAssertJSONStrictlyEqualsWhenJSONsAreEqual(t *testing.T) {
    tt := &MockTesting{}
    AssertJSONEquals(tt, `{"url": "URL 1"}`, `{"url": "URL 1"}`, true)
    assert.False(t, tt.Failed, tt.Message)
}

func TestAssertJSONStrictlyEqualsWhenJSONsAreNotEqual(t *testing.T) {
    tt := &MockTesting{}
    AssertJSONEquals(tt, `{"url": "URL 1"}`, `{"url": "URL 2"}`, true)
    assert.True(t, tt.Failed, tt.Message)
    assert.Equal(t, "[url:\nExpected: \"URL 1\"\ngot: \"URL 2\"\n]", tt.Message)
}

func TestAssertJSONStrictlyNotEqualsWhenJSONsAreEqual(t *testing.T) {
    tt := &MockTesting{}
    AssertJSONNotEquals(tt, `{"url": "URL 1"}`, `{"url": "URL 2"}`, true)
    assert.False(t, tt.Failed, tt.Message)
}

func TestAssertJSONStrictlyNotEqualsWhenJSONsAreNotEqual(t *testing.T) {
    tt := &MockTesting{}
    AssertJSONNotEquals(tt, `{"url": "URL 1"}`, `{"url": "URL 1"}`, true)
    assert.True(t, tt.Failed, tt.Message)
}
