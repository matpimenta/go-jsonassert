package jsonassert

import (
    "fmt"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
    testPass("\"Joe\"", "\"Joe\"", STRICT, t);
    testPass("\"Joe\"", "\"Joe\"", LENIENT, t);
    testPass("\"Joe\"", "\"Joe\"", NON_EXTENSIBLE, t);
    testPass("\"Joe\"", "\"Joe\"", STRICT_ORDER, t);
    testFail("\"Joe\"", "\"Joe1\"", STRICT, t);
    testFail("\"Joe\"", "\"Joe2\"", LENIENT, t);
    testFail("\"Joe\"", "\"Joe3\"", NON_EXTENSIBLE, t);
    testFail("\"Joe\"", "\"Joe4\"", STRICT_ORDER, t);
}

func TestNumber(t *testing.T) {
    testPass("123", "123", STRICT, t);
    testPass("123", "123", LENIENT, t);
    testPass("123", "123", NON_EXTENSIBLE, t);
    testPass("123", "123", STRICT_ORDER, t);
    testFail("123", "1231", STRICT, t);
    testFail("123", "1232", LENIENT, t);
    testFail("123", "1233", NON_EXTENSIBLE, t);
    testFail("123", "1234", STRICT_ORDER, t);
    testPass("0", "0", STRICT, t);
    testPass("-1", "-1", STRICT, t);
    testPass("0.1", "0.1", STRICT, t);
    testPass("1.2e5", "1.2e5", STRICT, t);
    testPass("20.4e-1", "20.4e-1", STRICT, t);
    testFail("310.1e-1", "31.01", STRICT, t); // should fail though numbers are the same?
}

func TestSimple(t *testing.T) {
    testPass(`{"id":1}`, `{"id":1}`, STRICT, t);
    testFail(`{"id":1}`, `{"id":2}`, STRICT, t);
    testPass(`{"id":1}`, `{"id":1}`, LENIENT, t);
    testFail(`{"id":1}`, `{"id":2}`, LENIENT, t);
    testPass(`{"id":1}`, `{"id":1}`, NON_EXTENSIBLE, t);
    testFail(`{"id":1}`, `{"id":2}`, NON_EXTENSIBLE, t);
    testPass(`{"id":1}`, `{"id":1}`, STRICT_ORDER, t);
    testFail(`{"id":1}`, `{"id":2}`, STRICT_ORDER, t);
}

func TestSimpleStrict(t *testing.T) {
    testPass(`{"id":1}`, `{"id":1,"name":"Joe"}`, LENIENT, t);
    testFail(`{"id":1}`, `{"id":1,"name":"Joe"}`, STRICT, t);
    testPass(`{"id":1}`, `{"id":1,"name":"Joe"}`, STRICT_ORDER, t);
    testFail(`{"id":1}`, `{"id":1,"name":"Joe"}`, NON_EXTENSIBLE, t);
}

func TestReversed(t *testing.T) {
    testPass(`{"name":"Joe","id":1}`, `{"id":1,"name":"Joe"}`, LENIENT, t);
    testPass(`{"name":"Joe","id":1}`, `{"id":1,"name":"Joe"}`, STRICT, t);
    testPass(`{"name":"Joe","id":1}`, `{"id":1,"name":"Joe"}`, NON_EXTENSIBLE, t);
    testPass(`{"name":"Joe","id":1}`, `{"id":1,"name":"Joe"}`, STRICT_ORDER, t);
}

func TestArray(t *testing.T) {
    testPass("[1,2,3]","[1,2,3]", STRICT, t);
    testPass("[1,2,3]","[1,3,2]", LENIENT, t);
    testFail("[1,2,3]","[1,3,2]", STRICT, t);
    testFail("[1,2,3]","[4,5,6]", LENIENT, t);
    testPass("[1,2,3]","[1,2,3]", STRICT_ORDER, t);
    testPass("[1,2,3]","[1,3,2]", NON_EXTENSIBLE, t);
    testFail("[1,2,3]","[1,3,2]", STRICT_ORDER, t);
    testFail("[1,2,3]","[4,5,6]", NON_EXTENSIBLE, t);
}

func TestNested(t *testing.T) {
    testPass(`{"id":1,"address":{"addr1":"123 Main", "addr2":null, "city":"Houston", "state":"TX"}}`,
            `{"id":1,"address":{"addr1":"123 Main", "addr2":null, "city":"Houston", "state":"TX"}}`, STRICT, t);
    testFail(`{"id":1,"address":{"addr1":"123 Main", "addr2":null, "city":"Houston", "state":"TX"}}`,
            `{"id":1,"address":{"addr1":"123 Main", "addr2":null, "city":"Austin", "state":"TX"}}`, STRICT, t);
}

func TestVeryNested(t *testing.T) {
    testPass(`{"a":{"b":{"c":{"d":{"e":{"f":{"g":{"h":{"i":{"j":{"k":{"l":{"m":{"n":{"o":{"p":"blah"}}}}}}}}}}}}}}}}`,
            `{"a":{"b":{"c":{"d":{"e":{"f":{"g":{"h":{"i":{"j":{"k":{"l":{"m":{"n":{"o":{"p":"blah"}}}}}}}}}}}}}}}}`, STRICT, t);
    testFail(`{"a":{"b":{"c":{"d":{"e":{"f":{"g":{"h":{"i":{"j":{"k":{"l":{"m":{"n":{"o":{"p":"blah"}}}}}}}}}}}}}}}}`,
            `{"a":{"b":{"c":{"d":{"e":{"f":{"g":{"h":{"i":{"j":{"k":{"l":{"m":{"n":{"o":{"z":"blah"}}}}}}}}}}}}}}}}`, STRICT, t);
}

func TestSimpleArray(t *testing.T) {
    testPass(`{"id":1,"pets":["dog","cat","fish"]}`, // Exact to exact (strict)
            `{"id":1,"pets":["dog","cat","fish"]}`,
            STRICT, t);
    testFail(`{"id":1,"pets":["dog","cat","fish"]}`, // Out-of-order fails (strict)
            `{"id":1,"pets":["dog","fish","cat"]}`,
            STRICT, t);
    testPass(`{"id":1,"pets":["dog","cat","fish"]}`, // Out-of-order ok
            `{"id":1,"pets":["dog","fish","cat"]}`,
            LENIENT, t);
    testPass(`{"id":1,"pets":["dog","cat","fish"]}`, // Out-of-order ok
            `{"id":1,"pets":["dog","fish","cat"]}`,
            NON_EXTENSIBLE, t);
    testFail(`{"id":1,"pets":["dog","cat","fish"]}`, // Out-of-order fails (strict order)
            `{"id":1,"pets":["dog","fish","cat"]}`,
            STRICT_ORDER, t);
    testFail(`{"id":1,"pets":["dog","cat","fish"]}`, // Mismatch
            `{"id":1,"pets":["dog","cat","bird"]}`,
            STRICT, t);
    testFail(`{"id":1,"pets":["dog","cat","fish"]}`, // Mismatch
            `{"id":1,"pets":["dog","cat","bird"]}`,
            LENIENT, t);
    testFail(`{"id":1,"pets":["dog","cat","fish"]}`, // Mismatch
            `{"id":1,"pets":["dog","cat","bird"]}`,
            STRICT_ORDER, t);
    testFail(`{"id":1,"pets":["dog","cat","fish"]}`, // Mismatch
            `{"id":1,"pets":["dog","cat","bird"]}`,
            NON_EXTENSIBLE, t);
}

func TestSimpleMixedArray(t *testing.T) {
    testPass(`{"stuff":[321, "abc"]}`, `{"stuff":["abc", 321]}`, LENIENT, t);
    testFail(`{"stuff":[321, "abc"]}`, `{"stuff":["abc", 789]}`, LENIENT, t);
}

func TestComplexMixedStrictArray(t *testing.T) {
    testPass(`{"stuff":[{"pet":"cat"},{"car":"Ford"}]}`, `{"stuff":[{"pet":"cat"},{"car":"Ford"}]}`, STRICT, t);
}

func TestComplexMixedArray(t *testing.T) {
    testPass(`{"stuff":[{"pet":"cat"},{"car":"Ford"}]}`, `{"stuff":[{"pet":"cat"},{"car":"Ford"}]}`, LENIENT, t);
}

func TestComplexArrayNoUniqueID(t *testing.T) {
    testPass(`{"stuff":[{"address":{"addr1":"123 Main"}}, {"address":{"addr1":"234 Broad"}}]}`,
            `{"stuff":[{"address":{"addr1":"123 Main"}}, {"address":{"addr1":"234 Broad"}}]}`,
            LENIENT, t);
}

func TestSimpleAndComplexStrictArray(t *testing.T) {
    testPass(`{"stuff":[123,{"a":"b"}]}`, `{"stuff":[123,{"a":"b"}]}`, STRICT, t);
}

func TestSimpleAndComplexArray(t *testing.T) {
    testPass(`{"stuff":[123,{"a":"b"}]}`, `{"stuff":[123,{"a":"b"}]}`, LENIENT, t);
}

func TestComplexArray(t *testing.T) {
    testPass(`{"id":1,"name":"Joe","friends":[{"id":2,"name":"Pat","pets":["dog"]},{"id":3,"name":"Sue","pets":["bird","fish"]}],"pets":[]}`,
             `{"id":1,"name":"Joe","friends":[{"id":2,"name":"Pat","pets":["dog"]},{"id":3,"name":"Sue","pets":["bird","fish"]}],"pets":[]}`,
             STRICT, t); // Exact to exact (strict)
    testFail(`{"id":1,"name":"Joe","friends":[{"id":2,"name":"Pat","pets":["dog"]},{"id":3,"name":"Sue","pets":["bird","fish"]}],"pets":[]}`,
            `{"id":1,"name":"Joe","friends":[{"id":3,"name":"Sue","pets":["fish","bird"]},{"id":2,"name":"Pat","pets":["dog"]}],"pets":[]}`,
            STRICT, t); // Out-of-order fails (strict)
    testFail(`{"id":1,"name":"Joe","friends":[{"id":2,"name":"Pat","pets":["dog"]},{"id":3,"name":"Sue","pets":["bird","fish"]}],"pets":[]}`,
            `{"id":1,"name":"Joe","friends":[{"id":3,"name":"Sue","pets":["fish","bird"]},{"id":2,"name":"Pat","pets":["dog"]}],"pets":[]}`,
            STRICT_ORDER, t); // Out-of-order fails (strict order)
    testPass(`{"id":1,"name":"Joe","friends":[{"id":2,"name":"Pat","pets":["dog"]},{"id":3,"name":"Sue","pets":["bird","fish"]}],"pets":[]}`,
            `{"id":1,"name":"Joe","friends":[{"id":3,"name":"Sue","pets":["fish","bird"]},{"id":2,"name":"Pat","pets":["dog"]}],"pets":[]}`,
            LENIENT, t); // Out-of-order ok
    testPass(`{"id":1,"name":"Joe","friends":[{"id":2,"name":"Pat","pets":["dog"]},{"id":3,"name":"Sue","pets":["bird","fish"]}],"pets":[]}`,
            `{"id":1,"name":"Joe","friends":[{"id":3,"name":"Sue","pets":["fish","bird"]},{"id":2,"name":"Pat","pets":["dog"]}],"pets":[]}`,
            NON_EXTENSIBLE, t); // Out-of-order ok
    testFail(`{"id":1,"name":"Joe","friends":[{"id":2,"name":"Pat","pets":["dog"]},{"id":3,"name":"Sue","pets":["bird","fish"]}],"pets":[]}`,
            `{"id":1,"name":"Joe","friends":[{"id":2,"name":"Pat","pets":["dog"]},{"id":3,"name":"Sue","pets":["cat","fish"]}],"pets":[]}`,
            STRICT, t); // Mismatch (strict)
    testFail(`{"id":1,"name":"Joe","friends":[{"id":2,"name":"Pat","pets":["dog"]},{"id":3,"name":"Sue","pets":["bird","fish"]}],"pets":[]}`,
            `{"id":1,"name":"Joe","friends":[{"id":2,"name":"Pat","pets":["dog"]},{"id":3,"name":"Sue","pets":["cat","fish"]}],"pets":[]}`,
            LENIENT, t); // Mismatch
    testFail(`{"id":1,"name":"Joe","friends":[{"id":2,"name":"Pat","pets":["dog"]},{"id":3,"name":"Sue","pets":["bird","fish"]}],"pets":[]}`,
            `{"id":1,"name":"Joe","friends":[{"id":2,"name":"Pat","pets":["dog"]},{"id":3,"name":"Sue","pets":["cat","fish"]}],"pets":[]}`,
            STRICT_ORDER, t); // Mismatch
    testFail(`{"id":1,"name":"Joe","friends":[{"id":2,"name":"Pat","pets":["dog"]},{"id":3,"name":"Sue","pets":["bird","fish"]}],"pets":[]}`,
            `{"id":1,"name":"Joe","friends":[{"id":2,"name":"Pat","pets":["dog"]},{"id":3,"name":"Sue","pets":["cat","fish"]}],"pets":[]}`,
            NON_EXTENSIBLE, t); // Mismatch
}

func TestArrayOfArraysStrict(t *testing.T) {
    testPass(`{"id":1,"stuff":[[1,2],[2,3],[],[3,4]]}`, `{"id":1,"stuff":[[1,2],[2,3],[],[3,4]]}`, STRICT, t);
    testFail(`{"id":1,"stuff":[[1,2],[2,3],[3,4],[]]}`, `{"id":1,"stuff":[[1,2],[2,3],[],[3,4]]}`, STRICT, t);
}

func TestArrayOfArrays(t *testing.T) {
    testPass(`{"id":1,"stuff":[[4,3],[3,2],[],[1,2]]}`, `{"id":1,"stuff":[[1,2],[2,3],[],[3,4]]}`, LENIENT, t);
}

func TestLenientArrayRecursion(t *testing.T) {
    testPass(`[{"arr":[5, 2, 1]}]`, `[{"b":3, "arr":[1, 5, 2]}]`, LENIENT, t);
}

func TestFieldMismatch(t *testing.T) {
    /*
    JSONCompareResult result = JSONCompare.compareJSON("{name:\"Pat\"}", "{name:\"Sue\"}", STRICT, t);
    FieldComparisonFailure comparisonFailure = result.getFieldFailures().iterator().next(, t);
    Assert.assertEquals("Pat", comparisonFailure.getExpected(), t);
    Assert.assertEquals("Sue", comparisonFailure.getActual(), t);
    Assert.assertEquals("name", comparisonFailure.getField(), t);
    */
}

func TestBooleanArray(t *testing.T) {
    testPass("[true, false, true, true, false]", "[true, false, true, true, false]", STRICT, t);
    testPass("[false, true, true, false, true]", "[true, false, true, true, false]", LENIENT, t);
    testFail("[false, true, true, false, true]", "[true, false, true, true, false]", STRICT, t);
    testPass("[false, true, true, false, true]", "[true, false, true, true, false]", NON_EXTENSIBLE, t);
    testFail("[false, true, true, false, true]", "[true, false, true, true, false]", STRICT_ORDER, t);
}

func TestNullProperty(t *testing.T) {
    testFail(`{"id":1,"name":"Joe"}`, `{"id":1,"name":null}`, STRICT, t);
    testFail(`{"id":1,"name":null}`, `{"id":1,"name":"Joe"}`, STRICT, t);
}

func TestIncorrectTypes(t *testing.T) {
    testFail(`{"id":1,"name":"Joe"}`, `{"id":1,"name":[]}`, STRICT, t);
    testFail(`{"id":1,"name":[]}`, `{"id":1,"name":"Joe"}`, STRICT, t);
}

func TestNullEquality(t *testing.T) {
    testPass(`{"id":1,"name":null}`, `{"id":1,"name":null}`, STRICT, t);
}

func TestExpectedArrayButActualObject(t *testing.T) {
    testFail("[1]", `{"id":1}`, LENIENT, t);
}

func TestExpectedObjectButActualArray(t *testing.T) {
    testFail(`{"id":1}`, "[1]", LENIENT, t);
}

func testPass(expected string, actual string, compareMode JSONCompareMode, t *testing.T) {
    message := fmt.Sprintf("%s == %s (%d)", expected, actual, compareMode)
    jsonCompare := NewJSONCompare()
    result := jsonCompare.CompareJSON(expected, actual, compareMode)
    assert.True(t, result.Passed(), message)
}

func testFail(expected string, actual string, compareMode JSONCompareMode, t *testing.T) {
    message := fmt.Sprintf("%s != %s (%d)", expected, actual, compareMode)
    jsonCompare := NewJSONCompare()
    result := jsonCompare.CompareJSON(expected, actual, compareMode)
    assert.True(t, result.Failed(), message)
}
