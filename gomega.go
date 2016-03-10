package jsonassert

import (
    "errors"
    "fmt"
)

type JSONAssertMatcher struct {
    compareMode JSONCompareMode
    expected interface{}
    results *JSONCompareResult
}

func (m *JSONAssertMatcher)	Match(actual interface{}) (success bool, err error) {
    if expectedStr, actualStr, ok := valuesAsString(m.expected, actual); ok {
        compare := NewJSONCompare()
        m.results = compare.CompareJSON(expectedStr, actualStr, m.compareMode)

        return m.results.Passed(), nil
    } else {
        return false, errors.New("Invalid json")
    }
}

func (m *JSONAssertMatcher)	FailureMessage(actual interface{}) (message string) {
    if m.results != nil {
        return m.results.GetMessage()
    } else {
        return "Invalid json"
    }
}

func (m *JSONAssertMatcher) NegatedFailureMessage(actual interface{}) (message string) {
    if m.results != nil {
        return m.results.GetMessage()
    } else {
        return "Invalid json"
    }
}

func MatchJSONStrictly(expected interface{}) *JSONAssertMatcher {
    return &JSONAssertMatcher{ expected: expected, compareMode: STRICT }
}

func MatchJSONLeniently(expected interface{}) *JSONAssertMatcher {
    return &JSONAssertMatcher{ expected: expected, compareMode: LENIENT }
}

func valuesAsString(expected, actual interface{}) (string, string, bool) {
    expectedStr, err1 := asString(expected)
    actualStr, err2 := asString(actual)
    if err1 && err2 {
        return expectedStr, actualStr, true
    } else {
        return "", "", false
    }
}

func asString(value interface{}) (string, bool) {
    if bytes, ok := value.([]byte); ok {
        return string(bytes), true
    }

    if str, ok := value.(string); ok {
        return str, true
    }

    if stringer, ok := value.(fmt.Stringer); ok {
        return stringer.String(), true
    }

    return "", false
}
