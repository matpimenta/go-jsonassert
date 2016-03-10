package jsonassert

import (
    "testing"
	. "github.com/onsi/gomega"
)

func TestMatchJSONLenientlyWhenJSONsAreEqual(t *testing.T) {
    RegisterTestingT(t)
    Expect(`{"url": "URL 1", "name": "John"}`).To(MatchJSONLeniently(`{"url": "URL 1"}`))
}

func TestMatchJSONLenientlyWhenJSONsAreNotEqual(t *testing.T) {
    RegisterTestingT(t)
    Expect(`{"url": "URL 1", "name": "John"}`).NotTo(MatchJSONLeniently(`{"url": "URL 2"}`))
}

func TestMatchJSONStrictlyWhenJSONsAreEqual(t *testing.T) {
    RegisterTestingT(t)
    Expect(`{"url": "URL 1"}`).To(MatchJSONStrictly(`{"url": "URL 1"}`))
}

func TestMatchJSONStrictlyWhenJSONsAreNotEqual(t *testing.T) {
    RegisterTestingT(t)
    Expect(`{"url": "URL 1", "name": "John"}`).NotTo(MatchJSONStrictly(`{"url": "URL 1"}`))
}

