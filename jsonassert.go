package jsonassert

import (
    "github.com/bitly/go-simplejson"
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
}

func (res *JSONCompareResult) Fail(prefix string, expected interface{}, actual interface{}) {
    res.success = false
}

func (res *JSONCompareResult) Unexpected(message string, object interface{}) {
    res.success = false
}
