# expec
Expectations testing package for Go

## Usage

Import it directly into test. Only exported function is `Expec()`

## Docs

Expec is largely inspired by RSpec expectations, and made as close as possible.
Docs are pending.

## Example

```go
import (
	"testing"
	_ "github.com/dmajkic/expec"
)

func TestMe (t *testing.T) {
	variable := "Some string"

  // Equality
	Expec(t, variable).To.Eql("Some string")
	Expec(t, variable).NotTo.Eql("Some otherstring")

	// Match regexp
	Expec(t, variable).To.Match("string$")

  // Nil testing
	Expec(t, variable).NotTo.BeNil()

	// Errors
	variable2, err := someCall()
	Expec(t, err).To.BeNil()

	// Errors expected
	Expec(t, err).NoTo.BeNil()
	Expec(t, err).To.RaiseErr()
	Expec(t, err).To.RaiseErr("Something went wrong")
	Expec(t, err).To.RaiseErr(os.ErrNotExist)

  // Interfaces
	Expec(t, errors.new("error") ).To.BeA(error(nil))

  // Slices
	Expec(t, []int{1,2,3,4,5}).To.Include(3)
}
```
