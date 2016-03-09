package jsonassert

import (
    "fmt"
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
    expected, errE := NewJSONNodeFromString(expectedStr)
    actual, errA := NewJSONNodeFromString(actualStr)

    if errE != nil {
        return &JSONCompareResult{success: false}
    }
    if errA != nil {
        return &JSONCompareResult{success: false}
    }
    comparator := DefaultComparator{compareMode: compareMode}
    if expected.IsArray() {
        if actual.IsArray() {
            return comparator.CompareJSONArray(expected, actual)
        } else {
            result := NewJSONCompareResult()
            result.Fail("", expected, actual)
            return result
        }
    } else if expected.IsMap() {
        if actual.IsMap() {
            return comparator.CompareJSONObject(expected, actual)
        } else {
            result := NewJSONCompareResult()
            result.Fail("", expected, actual)
            return result
        }
    } else {
        result := NewJSONCompareResult()
        comparator.CompareValues("", expected.GetData(), actual.GetData(), result)
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
        if safe, ok := item.(JSONNode); ok {
            out := safe.String()
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
