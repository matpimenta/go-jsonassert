package jsonassert

import (
    "testing"
	. "github.com/onsi/gomega"
)

func TestLenientlyMatchJSONWhenJSONsAreEqual(t *testing.T) {
    RegisterTestingT(t)
    Expect(`{"url": "URL 1", "name": "John"}`).To(LenientlyMatchJSON(`{"url": "URL 1"}`))
}

func TestLenientlyMatchJSONWhenJSONsAreNotEqual(t *testing.T) {
    RegisterTestingT(t)
    Expect(`{"url": "URL 1", "name": "John"}`).NotTo(LenientlyMatchJSON(`{"url": "URL 2"}`))
}

func TestStrictlyMatchJSONWhenJSONsAreEqual(t *testing.T) {
    RegisterTestingT(t)
    Expect(`{"url": "URL 1"}`).To(StrictlyMatchJSON(`{"url": "URL 1"}`))
}

func TestStrictlyMatchJSONWhenJSONsAreNotEqual(t *testing.T) {
    RegisterTestingT(t)
    Expect(`{"url": "URL 1", "name": "John"}`).NotTo(StrictlyMatchJSON(`{"url": "URL 1"}`))
}

