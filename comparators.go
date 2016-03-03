package jsonassert

import (
    "fmt"
    "github.com/bitly/go-simplejson"
    "reflect"
)

type JSONComparator interface {
    CompareJSONObject(expected *simplejson.Json, actual *simplejson.Json) *JSONCompareResult
    CompareJSONArray(expected *simplejson.Json, actual *simplejson.Json) *JSONCompareResult
    CompareJSONObjectWithPrefix(prefix string, expected *simplejson.Json, actual *simplejson.Json, result *JSONCompareResult)
    CompareJSONArrayWithPrefix(prefix string, expected *simplejson.Json, actual *simplejson.Json, result *JSONCompareResult)
    CompareValues(prefix string, expected interface{}, actual interface{}, result *JSONCompareResult)
}

type DefaultComparator struct {
    compareMode JSONCompareMode
}

func (comp *DefaultComparator) CompareJSONObjectWithPrefix(prefix string, expected *simplejson.Json, actual *simplejson.Json, result *JSONCompareResult) {
    comp.CheckJsonObjectKeysExpectedInActual(prefix, expected, actual, result)
    if comp.compareMode == NON_EXTENSIBLE || comp.compareMode == STRICT {
        comp.CheckJsonObjectKeysActualInExpected(prefix, expected, actual, result)
    }
}

func (comp *DefaultComparator) CompareJSONArrayWithPrefix(prefix string, expected *simplejson.Json, actual *simplejson.Json, result *JSONCompareResult) {
    if len(expected.MustArray()) != len(actual.MustArray()) {
        result.FailWithMessage(fmt.Sprintf("%s[]: Expected %d vallues but got %d", prefix, len(expected.MustArray()), len(actual.MustArray())))
    } else if len(expected.MustArray()) == 0 {
        return
    }

    if comp.compareMode == STRICT || comp.compareMode == STRICT_ORDER {
        comp.CompareJSONArrayWithStrictOrder(prefix, expected, actual, result)
    } else {
        comp.RecursivelyCompareJSONArray(prefix, expected, actual, result)
    }
}

func (comp *DefaultComparator) CompareValues(prefix string, expected interface{}, actual interface{}, result *JSONCompareResult) {
    if actual != nil && expected == nil || actual == nil && expected != nil {
        result.Fail(prefix, expected, actual)
    } else if reflect.TypeOf(actual).Name() != reflect.TypeOf(expected).Name() ||
    reflect.TypeOf(actual).Kind() != reflect.TypeOf(expected).Kind() {
        result.Fail(prefix, expected, actual)
    } else {
        if expectedElementSafe, ok := expected.(*simplejson.Json); ok {
            if actualElementSafe, ok := actual.(*simplejson.Json); ok {
                if _, err := expectedElementSafe.Map(); err == nil {
                    if _, err := actualElementSafe.Map(); err == nil {
                        comp.CompareJSONObjectWithPrefix(prefix, expectedElementSafe, actualElementSafe, result)
                    } else {
                        result.Fail(prefix, expected, actual)
                    }
                } else if _, err := expectedElementSafe.Array(); err == nil {
                    if _, err := actualElementSafe.Array(); err == nil {
                        comp.CompareJSONArrayWithPrefix(prefix, expectedElementSafe, actualElementSafe, result)
                    } else {
                        result.Fail(prefix, expected, actual)
                    }
                } else {
                    if expectedElementSafe.Interface() != actualElementSafe.Interface() {
                        result.Fail(prefix, expected, actual)
                    }
                }
            } else {
                result.Fail(prefix, expected, actual)
            }

        } else {
            if reflect.TypeOf(expected).Kind() == reflect.Slice {
                if reflect.TypeOf(actual).Kind() == reflect.Slice {
                    expectedElementSafe := expected.([]interface{})
                    actualElementSafe := actual.([]interface{})

                    newExpected := simplejson.New()
                    newExpected.SetPath([]string{}, expectedElementSafe)
                    newActual := simplejson.New()
                    newActual.SetPath([]string{}, actualElementSafe)
                    if comp.CompareJSONArray(newExpected, newActual).Failed() {
                        result.Fail(prefix, expected, actual)
                    }
                } else {
                    result.Fail(prefix, expected, actual)
                }
            } else if reflect.TypeOf(expected).Kind() == reflect.Map {
                if reflect.TypeOf(actual).Kind() == reflect.Map {
                    expectedElementSafe := expected.(map[string]interface{})
                    actualElementSafe := actual.(map[string]interface{})

                    newExpected := mapToJson(expectedElementSafe)
                    newActual := mapToJson(actualElementSafe)
                    if comp.CompareJSONObject(newExpected, newActual).Failed() {
                        result.Fail(prefix, expected, actual)
                    }
                } else {
                    result.Fail(prefix, expected, actual)
                }

            } else if expected != actual {
                result.Fail(prefix, expected, actual)
            }
        }
    }

}

func (comp *DefaultComparator) CompareJSONObject(expected *simplejson.Json, actual *simplejson.Json) *JSONCompareResult {
    result := NewJSONCompareResult()
    comp.CompareJSONObjectWithPrefix("", expected, actual, result)
    return result
}

func (comp *DefaultComparator) CompareJSONArray(expected *simplejson.Json, actual *simplejson.Json) *JSONCompareResult {
    result := NewJSONCompareResult()
    comp.CompareJSONArrayWithPrefix("", expected, actual, result)
    return result
}

func (comp *DefaultComparator) CompareJSONArrayOfSimpleValues(key string, expected simplejson.Json, actual simplejson.Json, result *JSONCompareResult) {
    expectedCount := getCardinalityMap(jsonArrayToList(expected))
    actualCount := getCardinalityMap(jsonArrayToList(expected))
    for id, expectedValue := range expectedCount {
        if actualValue, ok := actualCount[id]; ok {
            if actualValue != expectedValue {
                result.FailWithMessage(fmt.Sprintf("%s[]: Expected %d occurrences(s) of %s but got %d occurrence(s)", key, expectedValue, id, actualValue))
            }
        } else {
            result.Missing(fmt.Sprintf("%s[]", key), id)
        }
    }

    for id, _ := range actualCount {
        if _, ok := expectedCount[id]; ok {
            result.Unexpected(fmt.Sprintf("%s[]", key), id)
        }

    }
}

func (comp *DefaultComparator) CompareJSONArrayWithStrictOrder(key string, expected *simplejson.Json, actual *simplejson.Json, result *JSONCompareResult) {
    for i, expectedValue := range expected.MustArray() {
        actualValues := actual.MustArray()
        var actualValue interface{}
        if i < len(actualValues) {
            actualValue = actualValues[i]
        }
        comp.CompareValues(fmt.Sprintf("%s[%d]", key, i), expectedValue, actualValue, result)
    }
}

func (comp *DefaultComparator) RecursivelyCompareJSONArray(key string, expected *simplejson.Json, actual *simplejson.Json, result *JSONCompareResult) {
    matched := []int{}

    for i, expectedElement := range expected.MustArray() {
        matchFound := false
        for j, actualElement := range actual.MustArray() {
            if contains(matched, j) ||
            reflect.TypeOf(actualElement).Name() != reflect.TypeOf(expectedElement).Name() ||
            reflect.TypeOf(actualElement).Kind() != reflect.TypeOf(expectedElement).Kind() {
                continue
            }
            if expectedElementSafe, ok := expectedElement.(*simplejson.Json); ok {
                if actualElementSafe, ok := actualElement.(*simplejson.Json); ok {
                    if _, err := expectedElementSafe.Map(); err == nil {
                        if _, err := actualElementSafe.Map(); err == nil {
                            if comp.CompareJSONObject(expectedElementSafe, actualElementSafe).Passed() {
                                matched = append(matched, j)
                                matchFound = true
                                break
                            }
                        }
                    } else if _, err := expectedElementSafe.Array(); err == nil {
                        if _, err := actualElementSafe.Array(); err == nil {
                            if comp.CompareJSONArray(expectedElementSafe, actualElementSafe).Passed() {
                                matched = append(matched, j)
                                matchFound = true
                                break
                            }
                        }
                    }
                }
            } else {

                if reflect.TypeOf(expectedElement).Kind() == reflect.Slice {
                    if reflect.TypeOf(actualElement).Kind() == reflect.Slice {
                        expectedElementSafe := expectedElement.([]interface{})
                        actualElementSafe := actualElement.([]interface{})

                        newExpected := simplejson.New()
                        newExpected.SetPath([]string{}, expectedElementSafe)
                        newActual := simplejson.New()
                        newActual.SetPath([]string{}, actualElementSafe)
                        if comp.CompareJSONArray(newExpected, newActual).Passed() {
                            matched = append(matched, j)
                            matchFound = true
                            break
                        }
                    }
                } else if reflect.TypeOf(expectedElement).Kind() == reflect.Map {
                    if reflect.TypeOf(actualElement).Kind() == reflect.Map {
                        expectedElementSafe := expectedElement.(map[string]interface{})
                        actualElementSafe := actualElement.(map[string]interface{})

                        newExpected := mapToJson(expectedElementSafe)
                        newActual := mapToJson(actualElementSafe)
                        if comp.CompareJSONObject(newExpected, newActual).Passed() {
                            matched = append(matched, j)
                            matchFound = true
                            break
                        }
                    }
                } else {
                    if expectedElement == actualElement {
                        matched = append(matched, j)
                        matchFound = true
                        break
                    }
                }
            }

        }
        if !matchFound {
            result.FailWithMessage(fmt.Sprintf("%s[%d] Could not find match for element %s", key, i, expectedElement))
        }
    }
}

func (comp *DefaultComparator) CheckJsonObjectKeysActualInExpected(prefix string, expected *simplejson.Json, actual *simplejson.Json, result *JSONCompareResult) {
    for key, _ := range actual.MustMap() {
        if _, ok := expected.CheckGet(key); !ok {
            result.Unexpected(prefix, key)
        }
    }
}

func (comp *DefaultComparator) CheckJsonObjectKeysExpectedInActual(prefix string, expected *simplejson.Json, actual *simplejson.Json, result *JSONCompareResult) {
    for key, _ := range expected.MustMap() {
        expectedValue, _ := expected.CheckGet(key)
        if actualValue, ok := actual.CheckGet(key); ok {
            comp.CompareValues(qualify(prefix, key), expectedValue, actualValue, result)
        } else {
            result.Missing(prefix, key)
        }
    }
}

func contains(s []int, e int) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func mapToJson(mymap map[string]interface{}) *simplejson.Json {
    json := simplejson.New()
    for k, v := range mymap {
        json.Set(k, v)
    }
    return json
}

func jsonArrayToList(expected simplejson.Json) []interface{} {
    return expected.MustArray()
}

func qualify(prefix string, key string) string {
    if prefix == "" {
        return key
    } else {
        return fmt.Sprintf("%s.%s", prefix, key)
    }
}

func getCardinalityMap(coll []interface{}) map[interface{}]int {
    count := map[interface{}]int{}
    for _, item := range coll {
        if _, ok := count[item]; ok {
            count[item]++
        } else {
            count[item] = 1
        }
    }
    return count
}
