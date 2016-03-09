package jsonassert

type Testing interface {
	Errorf(format string, args ...interface{})
}

func AssertJSONEquals(t Testing, expected, actual string, strict bool) {
    jsonCompare := NewJSONCompare()
    var compareMode JSONCompareMode
    if strict {
        compareMode = STRICT
    } else {
        compareMode = LENIENT
    }
    result := jsonCompare.CompareJSON(expected, actual, compareMode)
    if result.Failed() {
        t.Errorf("%s", result.GetMessage())
    }
}

func AssertJSONNotEquals(t Testing, expected, actual string, strict bool) {
    jsonCompare := NewJSONCompare()
    var compareMode JSONCompareMode
    if strict {
        compareMode = STRICT
    } else {
        compareMode = LENIENT
    }
    result := jsonCompare.CompareJSON(expected, actual, compareMode)
    if !result.Failed() {
        t.Errorf("Json comparison failed:\n%s\n", result.GetMessage())
    }

}

