package jsonassert

import (
    "fmt"
    "testing"
)

type MockTesting struct {
    Failed bool
    Message string
}

func (t *MockTesting) Errorf(format string, args ...interface{}) {
    t.Failed = true
    t.Message = fmt.Sprintf(format, args)
}

func TestAssertEqualsWhenJSONsAreEquals(t *testing.T) {
    tt := &MockTesting{}
    assert := NewJSONAssertions(tt)
    assert.AssertEquals(`{"url": "URL 1"}`, `{"url": "URL 1"}`, false)
    if tt.Failed {
        t.Errorf(tt.Message)
    }
}

func TestAssertEqualsWhenJSONsAreNotEquals(t *testing.T) {
    tt := &MockTesting{}
    assert := NewJSONAssertions(tt)
    assert.AssertEquals(`{"url": "URL 1"}`, `{"url": "URL 2"}`, false)
    if !tt.Failed {
        t.Errorf(tt.Message)
    } else {
        expected := "[url:\nExpected: \"URL 1\"\ngot: \"URL 2\"\n]"
        if tt.Message != expected {
            t.Errorf("Invalid message: %s", tt.Message)
        }
    }
}
