JSONAssert for GO
=================

Write JSON unit tests in less code. Great for testing REST interfaces.

This is a port from the Java library https://github.com/skyscreamer/JSONassert

Usage
-----

```go
import (
	_ "github.com/matpimenta/go-jsonassert"
)

func TestExample(t *testing.T) {
    actualResult := responseFromRestEndpoint()
    AssertJSONEquals(t, `{ "url": "Repository 1" }`, actualResult, false)
}

```
