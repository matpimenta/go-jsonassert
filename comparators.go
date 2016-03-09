package jsonassert

import (
	"fmt"
	"reflect"
)

type JSONComparator interface {
	CompareJSONObject(expected, actual JSONNode) *JSONCompareResult
	CompareJSONArray(expected, actual JSONNode) *JSONCompareResult
	CompareJSONObjectWithPrefix(prefix string, expected, actual JSONNode, result *JSONCompareResult)
	CompareJSONArrayWithPrefix(prefix string, expected, actual JSONNode, result *JSONCompareResult)
	CompareValues(prefix string, expected interface{}, actual interface{}, result *JSONCompareResult)
}

type DefaultComparator struct {
	compareMode JSONCompareMode
}

func (comp *DefaultComparator) CompareJSONObjectWithPrefix(prefix string, expected, actual JSONNode, result *JSONCompareResult) {
	comp.CheckJsonObjectKeysExpectedInActual(prefix, expected, actual, result)
	if comp.compareMode == NON_EXTENSIBLE || comp.compareMode == STRICT {
		comp.CheckJsonObjectKeysActualInExpected(prefix, expected, actual, result)
	}
}

func (comp *DefaultComparator) CompareJSONArrayWithPrefix(prefix string, expected, actual JSONNode, result *JSONCompareResult) {
	if expected.GetSize() != actual.GetSize() {
		result.FailWithMessage(fmt.Sprintf("%s[]: Expected %d vallues but got %d", prefix, expected.GetSize(), actual.GetSize()))
	} else if expected.GetSize() == 0 {
		return
	}

	if comp.compareMode == STRICT || comp.compareMode == STRICT_ORDER {
		comp.CompareJSONArrayWithStrictOrder(prefix, expected, actual, result)
	} else {
		comp.RecursivelyCompareJSONArray(prefix, expected, actual, result)
	}
}

func (comp *DefaultComparator) CompareValues(prefix string, expected, actual interface{}, result *JSONCompareResult) {
	if actual != nil && expected == nil || actual == nil && expected != nil {
		result.Fail(prefix, expected, actual)
	} else if reflect.TypeOf(actual).Name() != reflect.TypeOf(expected).Name() ||
		reflect.TypeOf(actual).Kind() != reflect.TypeOf(expected).Kind() {

		result.Fail(prefix, expected, actual)
	} else {
		if expectedElementSafe, ok := expected.(JSONNode); ok {
			if actualElementSafe, ok := actual.(JSONNode); ok {
				comp.CompareValuesJSONNode(prefix, expectedElementSafe, actualElementSafe, result)
			} else {
				result.Fail(prefix, expected, actual)
			}

		} else {
			if reflect.TypeOf(expected).Kind() == reflect.Slice {
				if reflect.TypeOf(actual).Kind() == reflect.Slice {
					expectedElementSafe := expected.([]interface{})
					actualElementSafe := actual.([]interface{})

					newExpected := NewJSONNode()
					newExpected.SetArray(expectedElementSafe)
					newActual := NewJSONNode()
					newActual.SetArray(actualElementSafe)
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

					newExpected := NewJSONNodeFromMap(expectedElementSafe)
					newActual := NewJSONNodeFromMap(actualElementSafe)
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

func (comp *DefaultComparator) CompareValuesJSONNode(prefix string, expected, actual JSONNode, result *JSONCompareResult) {
	if expected.IsMap() {
		if actual.IsMap() {
			comp.CompareJSONObjectWithPrefix(prefix, expected, actual, result)
		} else {
			result.Fail(prefix, expected, actual)
		}
	} else if expected.IsArray() {
		if actual.IsArray() {
			comp.CompareJSONArrayWithPrefix(prefix, expected, actual, result)
		} else {
			result.Fail(prefix, expected, actual)
		}
	} else {
		if expected.GetData() != actual.GetData() {
			result.Fail(prefix, expected, actual)
		}
	}
}

func (comp *DefaultComparator) CompareJSONObject(expected, actual JSONNode) *JSONCompareResult {
	result := NewJSONCompareResult()
	comp.CompareJSONObjectWithPrefix("", expected, actual, result)
	return result
}

func (comp *DefaultComparator) CompareJSONArray(expected, actual JSONNode) *JSONCompareResult {
	result := NewJSONCompareResult()
	comp.CompareJSONArrayWithPrefix("", expected, actual, result)
	return result
}

func (comp *DefaultComparator) CompareJSONArrayWithStrictOrder(key string, expected, actual JSONNode, result *JSONCompareResult) {
	for i, expectedValue := range expected.GetArray() {
		actualValues := actual.GetArray()
		var actualValue interface{}
		if i < len(actualValues) {
			actualValue = actualValues[i]
		}
		comp.CompareValues(fmt.Sprintf("%s[%d]", key, i), expectedValue, actualValue, result)
	}
}

func (comp *DefaultComparator) RecursivelyCompareJSONArray(key string, expected, actual JSONNode, result *JSONCompareResult) {
	matched := []int{}

	for i, expectedElement := range expected.GetArray() {
		matchFound := false
		for j, actualElement := range actual.GetArray() {
			if contains(matched, j) ||
				reflect.TypeOf(actualElement).Name() != reflect.TypeOf(expectedElement).Name() ||
				reflect.TypeOf(actualElement).Kind() != reflect.TypeOf(expectedElement).Kind() {
				continue
			}
			if expectedElementSafe, ok := expectedElement.(JSONNode); ok {
				if actualElementSafe, ok := actualElement.(JSONNode); ok {
					if expectedElementSafe.IsMap() {
						if actualElementSafe.IsMap() {
							if comp.CompareJSONObject(expectedElementSafe, actualElementSafe).Passed() {
								matched = append(matched, j)
								matchFound = true
								break
							}
						}
					} else if expectedElementSafe.IsArray() {
						if actualElementSafe.IsArray() {
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

						newExpected := NewJSONNodeFromArray(expectedElementSafe)
						newActual := NewJSONNodeFromArray(actualElementSafe)
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

						newExpected := NewJSONNodeFromMap(expectedElementSafe)
						newActual := NewJSONNodeFromMap(actualElementSafe)
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

func (comp *DefaultComparator) CheckJsonObjectKeysActualInExpected(prefix string, expected, actual JSONNode, result *JSONCompareResult) {
	for key, _ := range actual.GetMap() {
		if _, ok := expected.CheckGet(key); !ok {
			result.Unexpected(prefix, key)
		}
	}
}

func (comp *DefaultComparator) CheckJsonObjectKeysExpectedInActual(prefix string, expected, actual JSONNode, result *JSONCompareResult) {
	for key, _ := range expected.GetMap() {
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

func qualify(prefix string, key string) string {
	if prefix == "" {
		return key
	} else {
		return fmt.Sprintf("%s.%s", prefix, key)
	}
}

