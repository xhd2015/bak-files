## Expected

- No harness error.
- `Response.Names` equals `["tmp", "binary", "vendor"]` in that order.

## Side Effects

- None.

## Errors

- Wrong order or missing names fails.

## Exit Code

- N/A

```go
import (
	"reflect"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	want := []string{"tmp", "binary", "vendor"}
	if !reflect.DeepEqual(resp.Names, want) {
		t.Fatalf("Flag.Names(): got %#v, want %#v", resp.Names, want)
	}
}
```
