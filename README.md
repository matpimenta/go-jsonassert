JSONAssert for GO
=================

Write JSON unit tests in less code. Great for testing REST interfaces.

This is a port from the Java library https://github.com/skyscreamer/JSONassert

Usage
-----

```go
package test

import (
    "testing"
	. "github.com/matpimenta/go-jsonassert"
)

func TestLenientAssertion(t *testing.T) {
    AssertJSONEquals(t, `{ "name": "John" }`, `{ "name": "John", "surname": "Smith" }`, false)
}

func TestStrictAssertion(t *testing.T) {
    AssertJSONEquals(t, `{ "name": "John" }`, `{ "name": "John" }`, true)
}

```

Gomega Matcher
--------------

```go
package jsonassert

import (
    "testing"
	. "github.com/onsi/gomega"
)

func TestMatchJSONLeniently(t *testing.T) {
    RegisterTestingT(t)
    Expect(`{"url": "URL 1", "name": "John"}`).To(MatchJSONLeniently(`{"url": "URL 1"}`))
}

func TestMatchJSONStrictly(t *testing.T) {
    RegisterTestingT(t)
    Expect(`{"url": "URL 1"}`).To(MatchJSONStrictly(`{"url": "URL 1"}`))
}

```

