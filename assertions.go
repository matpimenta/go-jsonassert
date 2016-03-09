package jsonassert

type Testing interface {
	Errorf(format string, args ...interface{})
}

type JSONAssertions struct {
    t Testing
}

func NewJSONAssertions(t Testing) *JSONAssertions {
    return &JSONAssertions{t: t}
}

func (a *JSONAssertions) AssertEquals(expected string, actual string, strict bool) {
    jsonCompare := NewJSONCompare()
    var compareMode JSONCompareMode
    if strict {
        compareMode = STRICT
    } else {
        compareMode = LENIENT
    }
    result := jsonCompare.CompareJSON(expected, actual, compareMode)
    if result.Failed() {
        a.t.Errorf("%s", result.GetMessage())
    }
}

func (a *JSONAssertions) AssertNotEquals(expected string, actual string, strict bool) {
    jsonCompare := NewJSONCompare()
    var compareMode JSONCompareMode
    if strict {
        compareMode = STRICT
    } else {
        compareMode = LENIENT
    }
    result := jsonCompare.CompareJSON(expected, actual, compareMode)
    if !result.Failed() {
        a.t.Errorf("Json comparison failed:\n%s\n", result.GetMessage())
    }

}

