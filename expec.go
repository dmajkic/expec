package expec

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

type Challenge struct {
	Subject
	not bool
}

type Subject struct {
	t *testing.T
	v interface{}
	To *Challenge
	NotTo *Challenge
	And *Challenge
}

func Expec(t *testing.T, v interface{}) *Subject {
	t.Helper()
	ret := &Subject{t:t, v:v}
	ret.To = &Challenge{*ret, false}
	ret.And = &Challenge{*ret, false}
	ret.NotTo = &Challenge{*ret, true}

	return ret
}

// expect is simple testing function which raises error if condition is not met
func (c *Challenge) expect(condition bool, eformat string, args ...interface{}) bool {
	c.t.Helper()

	if c.not {
		condition = !condition
	}

	if !condition {
		c.t.Errorf(eformat, args...)
	}
	return condition
}

// Eq expects both arguments to be equal, compared by ==
func (c *Challenge) Eq(expected interface{}) *Challenge {
	c.t.Helper()
	c.expect(expected == c.v, "Expected '%v' to equal '%v'", c.v, expected)
	return c
}

// Eql expects both arguments to be equal uses reflection.DeepEqual to perform test
func (c *Challenge) Eql(expected interface{}) *Challenge {
	c.t.Helper()
	c.expect(reflect.DeepEqual(c.v, expected), "Expected '%v' to equal '%v'", c.v, expected)
	return c
}

// Be expects both arguments to be equal uses reflection.DeepEqual to perform test
func (c *Challenge) Be(expected interface{}) *Challenge {
	c.t.Helper()
	c.expect(reflect.DeepEqual(c.v, expected), "Expected '%v' to equal '%v'", c.v, expected)
	return c
}

// BeNil should be used to check returned Go error.
func (c *Challenge) BeNil() *Challenge {
	c.t.Helper()
	v := reflect.ValueOf(c.v)
	c.expect(c.v == nil || v.IsZero() || v.IsNil(), "Expected '%v' to be nil", c.v)
	return c
}

// BeTrue should be used to check returned Go error.
func (c *Challenge) BeTrue() *Challenge {
	c.t.Helper()
	c.expect(c.v == true, "Expected '%v' to be true", c.v)
	return c
}

// BeTrue should be used to check returned Go error.
func (c *Challenge) BeFalse() *Challenge {
	c.t.Helper()
	c.expect(c.v == false, "Expected '%v' to be false", c.v)
	return c
}

// BeTrue should be used to check returned Go error.
func (c *Challenge) BeFalsy() *Challenge {
	c.t.Helper()
	c.expect(c.v == false || c.v == nil, "Expected '%v' to be false or nil", c.v)
	return c
}

// BeTrue should be used to check returned Go error.
func (c *Challenge) BeTruthy() *Challenge {
	c.t.Helper()
	c.expect(c.v != false && c.v != nil, "Expected '%v' to be false or nil", c.v)
	return c
}

// BeTrue should be used to check returned Go error.
func (c *Challenge) Match(pattern string) *Challenge {
	c.t.Helper()

	str, ok := c.v.(string)
	if !ok {
		if stringer, ok := c.v.(fmt.Stringer); ok {
			str = stringer.String()
		} else {
			c.t.Errorf("Error matching '%v' to '%v': not a string or Stringer", c.v, pattern)
			return c
		}
	}

	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		c.t.Errorf("Error matching '%v' to '%v': %v", c.v, pattern, err)
		return c
	}

	c.expect(matched, "Expected '%v' to match '%v'", c.v, pattern)
	return c
}

// BeA checks if subject is of same type or implements interface
func (c *Challenge) BeA(intf interface{}) *Challenge {
	c.t.Helper()
	t := reflect.TypeOf(c.v)
	intft := reflect.TypeOf(intf)

	c.expect((t == intft) || (intft.Kind() == reflect.Interface && t.Implements(intft)), "Expected '%v' to be a %v", c.v, intf)
	return c
}

// BeAn is an alias for BeA
func (c *Challenge) BeAn(intf interface{}) *Challenge {
	c.t.Helper()
	c.BeA(intf)
	return c
}

// Implement is an alias for BeA
func (c *Challenge) Implement(intf interface{}) *Challenge {
	c.t.Helper()
	c.BeA(intf)
	return c
}


// RaiseError expects Error to be raised
func (c *Challenge) RaiseError(arguments ...interface{}) *Challenge {
	c.t.Helper()
	if c.v != nil {
		c.expect(c.v != nil, "expected error to be raised, got nil")
		return c
	}

	err, ok := c.v.(error)
	if !ok {
		c.t.Errorf("Expected subject to be error, got %v", reflect.TypeOf(err))
		return c
	}

	switch len(arguments) {
	case 0:
		return c
	case 1:
		switch v := arguments[0].(type) {
		case error:
			c.expect(errors.Is(err, v), "expected error %v got %v", v, c.v)
		case string:
			c.expect(err.Error() == v, "expected error message:\n%v\n%v\n", v, err.Error())
		default:
			c.t.Errorf("Argument 1 should be error or string; got %v", arguments[0])
			return c
		}
	case 2:
		if target, ok := arguments[0].(error); !ok {
			c.t.Errorf("Argument 2 should be error; got %v", arguments[0])
			return c
		} else {
			c.expect(errors.Is(err, target), "expected error %v got %v", target, c.v)
		}

		if msg, ok := arguments[1].(string); !ok {
			c.t.Errorf("Argument 2 should be string; got %v", arguments[1])
			return c
		} else {
			c.expect(err.Error() == msg, "expected error message:\n%v\n%v\n", msg, err.Error())
		}
	default:
		c.t.Errorf("Expected no arguments, (error), (string) or (error, string)")
		return c
	}
	return c
}

func getSlice(v interface{}) (reflect.Value, error)  {
	a := reflect.ValueOf(v)
	if a.Kind() != reflect.Slice || a.Kind() != reflect.Array || a.Kind() != reflect.String  {
		return a, fmt.Errorf("array, string or slice expected, got %v", a.Type().Kind())
	}
	return a, nil
}

// Include checks that expected is included in array
func (c *Challenge) Include(expected interface{}) *Challenge {
	c.t.Helper()

	a, err := getSlice(expected)
	if err != nil {
		c.t.Error(err)
		return c
	}

	for i := 0; i < a.Len(); i++ {
		if reflect.DeepEqual(a.Index(i), expected) {
			return c
		}
	}

	c.t.Errorf("item %v is not incluuded in %v", expected, c.v)
	return c
}

