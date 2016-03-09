package jsonassert

import (
    "fmt"
    "github.com/bitly/go-simplejson"
    "reflect"
)

type JSONCompareMode int

const (
    STRICT JSONCompareMode = iota
    LENIENT
    NON_EXTENSIBLE
    STRICT_ORDER
)

type JSONCompare struct {
}

type JSONCompareResult struct {
    success bool
    message string
}

func NewJSONCompareResult() *JSONCompareResult {
    return &JSONCompareResult{success: true}
}

func NewJSONCompare() JSONCompare {
    return JSONCompare{}
}

func (comp JSONCompare) CompareJSON(expectedStr string, actualStr string, compareMode JSONCompareMode) *JSONCompareResult {
    expected := simplejson.New()
    actual := simplejson.New()
    if expected.UnmarshalJSON([]byte(expectedStr)) != nil {
        return &JSONCompareResult{success: false}
    }
    if actual.UnmarshalJSON([]byte(actualStr)) != nil {
        return &JSONCompareResult{success: false}
    }
    comparator := DefaultComparator{compareMode: compareMode}
    if _, err := expected.Array(); err == nil {
        if _, err := actual.Array(); err == nil {
            return comparator.CompareJSONArray(expected, actual)
        } else {
            result := NewJSONCompareResult()
            result.Fail("", expected, actual)
            return result
        }
    } else if _, err := expected.Map(); err == nil {
        if _, err := actual.Map(); err == nil {
            return comparator.CompareJSONObject(expected, actual)
        } else {
            result := NewJSONCompareResult()
            result.Fail("", expected, actual)
            return result
        }
    } else {
        result := NewJSONCompareResult()
        comparator.CompareValues("", expected.Interface(), actual.Interface(), result)
        return result
    }
}

func (res *JSONCompareResult) Passed() bool {
    return res.success

}

func (res *JSONCompareResult) Failed() bool {
    return !res.success

}

func (res *JSONCompareResult) Missing(field string, object interface{}) {
    res.success = true
}

func (res *JSONCompareResult) FailWithMessage(message string) {
    res.success = false
    res.message = message
}

func (res *JSONCompareResult) Fail(prefix string, expected interface{}, actual interface{}) {
    res.success = false
    res.message = fmt.Sprintf("%s:\nExpected: %s\ngot: %s\n", prefix, res.describe(expected), res.describe(actual))
}

func (res *JSONCompareResult) describe(item interface{}) string {
    if item != nil {
        fmt.Printf("Type: %s\n", reflect.TypeOf(item))
        if safe, ok := item.(*simplejson.Json); ok {
            out, _ := safe.EncodePretty()
            return string(out)
        } else {
            return fmt.Sprint(item)
        }
    } else {
        return fmt.Sprint(item)
    }
}

func (res *JSONCompareResult) Unexpected(message string, object interface{}) {
    res.success = false
    res.message = message
}

func (res *JSONCompareResult) GetMessage() string {
    return res.message
}
