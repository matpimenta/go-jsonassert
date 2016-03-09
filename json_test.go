package jsonassert

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJSONNodeArray(t *testing.T) {
	json, _ := NewJSONNodeFromString(`["a","b"]`)

	assert.True(t, json.IsArray(), "Should be an array")
	assert.False(t, json.IsMap(), "Should not be a map")
	assert.Equal(t, 2, json.GetSize())
}

func TestJSONNodeArrayWhenSetProgramatically(t *testing.T) {
	array := make([]interface{}, 2)
	array[0] = "a"
	array[1] = "b"
	json := NewJSONNodeFromArray(array)

	assert.True(t, json.IsArray(), "Should be an array")
	assert.False(t, json.IsMap(), "Should not be a map")
    assert.Equal(t, "a", json.GetArray()[0])
	assert.Equal(t, 2, json.GetSize())
}

func TestJSONNodeCheckGetForArrays(t *testing.T) {
	json, _ := NewJSONNodeFromString(`{"id":1,"array":["a","b"]}`)
	array, found := json.CheckGet("array")

	assert.True(t, found, "Should be found")
	assert.True(t, array.IsArray(), "Should be an array")
	assert.False(t, array.IsMap(), "Should not be a map")
	assert.Equal(t, 2, array.GetSize())
}

func TestJSONNodeCheckGetForString(t *testing.T) {
	json, _ := NewJSONNodeFromString(`{"id":1,"string":"a"}`)
	testString, found := json.CheckGet("string")

	assert.True(t, found, "Should be found")
	assert.False(t, testString.IsArray(), "Should not be an array")
	assert.False(t, testString.IsMap(), "Should not be a map")
	assert.Equal(t, 0, testString.GetSize())
	assert.Equal(t, "a", testString.GetData())
}
