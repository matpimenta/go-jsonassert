package jsonassert

import (
    "fmt"
    "github.com/bitly/go-simplejson"
)

type JSONNode interface {
    IsMap() bool
    IsArray() bool
    GetData() interface{}
    SetData(data interface{})
    SetArray(array []interface{})
    GetArray() []interface{}
    GetSize() int
    GetMap() map[string]interface{}
    CheckGet(key string) (JSONNode, bool)
    Set(key string, value interface{})
    String() string
}

type SimpleJSONNode struct {
    json *simplejson.Json
}

func NewJSONNode() JSONNode {
    return &SimpleJSONNode{
        json: simplejson.New(),
    }
}

func NewJSONNodeFromMap(mymap map[string]interface{}) JSONNode {
    json := NewJSONNode()
    for k, v := range mymap {
        json.Set(k, v)
    }
    return json
}

func NewJSONNodeFromArray(array []interface{}) JSONNode {
    json := NewJSONNode()
    json.SetArray(array)
    return json
}

func NewJSONNodeFromString(jsonStr string) (JSONNode, error) {
    json, err := simplejson.NewJson([]byte(jsonStr))
    if err == nil {
        return &SimpleJSONNode{ json: json }, nil
    }
    return nil, err
}

func (s *SimpleJSONNode) SetArray(array []interface{}) {
    s.json.SetPath([]string{}, array)
}

func (s *SimpleJSONNode) GetArray() []interface{} {
    return s.json.MustArray()
}

func (s *SimpleJSONNode) IsMap() bool {
    _, err := s.json.Map()
    if err == nil {
        return true
    }
    return false
}

func (s *SimpleJSONNode) IsArray() bool {
    _, err := s.json.Array()
    if err == nil {
        return true
    }
    return false
}

func (s *SimpleJSONNode) SetData(data interface{}) {
    s.json.SetPath([]string{}, data)
}

func (s *SimpleJSONNode) GetData() interface{} {
    return s.json.Interface()
}

func (s *SimpleJSONNode) GetSize() int {
    if s.IsArray() {
        return len(s.json.MustArray())
    } else if s.IsMap() {
        return len(s.json.MustMap())
    }
    return 0
}

func (s *SimpleJSONNode) GetMap() map[string]interface{} {
    return s.json.MustMap()
}

func (s *SimpleJSONNode) CheckGet(key string) (JSONNode, bool) {
    value, found := s.json.CheckGet(key)
    if found {
        json := &SimpleJSONNode{json: value}
        return json, true
    }
    return nil, false
}

func (s *SimpleJSONNode) Set(key string, value interface{}) {
    s.json.Set(key, value)
}

func (s *SimpleJSONNode) String() string {
    out, err := s.json.Encode()
    if err == nil {
        return string(out)
    }
    return fmt.Sprintf("%s", s.json)
}
