package expec

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// Subject is test subject that provides access to challenge its value
type Subject struct {
	t *testing.T
	v interface{}
	To *Challenge
	NotTo *Challenge
	And *Challenge
	Must *Challenge
}

// Challenge represents expetation subtest value
type Challenge struct {
	*Subject
	not bool
	fail bool
}

// Composing shortcuts
const (
	Eq = iota
	Eql
	Gt
	Lt
	Match
	Implement
	A
	Include
	StartWith
	EndWith
)


// Expec captures subject to challenge its value
func Expec(t *testing.T, v interface{}) *Subject {
	t.Helper()
	ret := &Subject{t:t, v:v}
	ret.To = &Challenge{ret, false, false}
	ret.And = &Challenge{ret, false, false}
	ret.NotTo = &Challenge{ret, true, false}
	ret.Must = &Challenge{ret, false, true}

	return ret
}

// expect is simple testing function which raises error if condition is not met
func (c *Challenge) expect(condition bool, eformat string, args ...interface{}) bool {
	c.t.Helper()

	if c.not {
		condition = !condition
		eformat = strings.ReplaceAll(eformat, " to ", " not to ")
	}

	if !condition {
		if c.fail {
			c.t.Fatalf(eformat, args...)
		} else {
			c.t.Errorf(eformat, args...)
		}
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
	if _, ok := c.v.(error); ok {
		return c.RaiseError(expected)
	}

	c.expect(reflect.DeepEqual(c.v, expected), "Expected '%v' to equal '%v'", c.v, expected)
	return c
}

// BeNil should be used to check returned Go error.
func (c *Challenge) BeNil() *Challenge {
	c.t.Helper()
	v := reflect.ValueOf(c.v)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		c.expect(c.v == nil || v.IsNil(), "Expected '%v' to be nil", c.v)
	}
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
			c.t.Fatalf("Error matching '%v' to '%v': not a string or Stringer", c.v, pattern)
			return c
		}
	}

	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		c.t.Fatalf("Error matching '%v' to '%v': %v", c.v, pattern, err)
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

	if c.v == nil {
		c.expect(false, "Value is nil. It does not implement %v", reflect.TypeOf(c.v))
		return c
	}

	t := reflect.TypeOf(c.v)

	intft := reflect.TypeOf(intf)
	if intft.Elem().Kind() != reflect.Interface {
		c.expect(false, "Expected '%v' to be pointer to interface like (*error)(nil)", intft)
		return c
	}
	c.expect(t == intft || (t.Implements(intft.Elem())), "Expected '%v' to be implementator of %v", t, intf)
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
		c.t.Fatalf("Expected subject to be error, got %v", reflect.TypeOf(err))
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
			c.t.Fatalf("Argument 1 should be error or string; got %v", arguments[0])
			return c
		}
	case 2:
		if target, ok := arguments[0].(error); !ok {
			c.t.Fatalf("Argument 2 should be error; got %v", arguments[0])
			return c
		} else {
			c.expect(errors.Is(err, target), "expected error %v got %v", target, c.v)
		}

		if msg, ok := arguments[1].(string); !ok {
			c.t.Fatalf("Argument 2 should be string; got %v", arguments[1])
			return c
		} else {
			c.expect(err.Error() == msg, "expected error message:\n%v\n%v\n", msg, err.Error())
		}
	default:
		c.t.Fatalf("Expected no arguments, (error), (string) or (error, string)")
		return c
	}
	return c
}

func getSlice(v interface{}) (reflect.Value, error)  {
	a := reflect.ValueOf(v)
	if a.Kind() != reflect.Slice && a.Kind() != reflect.Array && a.Kind() != reflect.String  {
		return a, fmt.Errorf("array, string or slice expected, got %v", a.Type().Kind())
	}
	return a, nil
}

func getString(values ...interface{}) (string, error)  {
	if len(values) == 0 {
		return "", nil
	}

	ret := ""

	for _, v := range values {
		s, ok := v.(string)
		if !ok {
			return ret, errors.New("string expected")
		}
		ret += s
	}

	return ret, nil
}

// Include checks that expected element is in array
func (c *Challenge) Include(elements ...interface{}) *Challenge {
	c.t.Helper()

	a, err := getSlice(c.v)
	if err != nil {
		c.t.Fatal(err)
		return c
	}

	if a.Kind() == reflect.String {
		s, err := getString(elements...)
		if err != nil {
			c.t.Fatal("argument must be string or array of strings")
			return c
		}
		c.expect(strings.Contains(a.String(), s), "expected '%v' to include '%v'", c.v, s)
		return c
	}

	found := make(map[interface{}]struct{}, len(elements))

	for i := 0; i < a.Len(); i++ {
		intf := a.Index(i).Interface()

		for _, e := range elements {
			if reflect.DeepEqual(intf, e) {
				found[intf] = struct{}{}
				if len(found) == len(elements) {
					break
				}
			}
		}

		if len(found) == len(elements) {
			break
		}
	}

	c.expect(len(found) == len(elements),  "expected '%v' to include '%v'", c.v, elements)
	return c
}

// StartWith checks that subject starts with provided elements
func (c *Challenge) StartWith(values ...interface{}) *Challenge {
	c.t.Helper()

	a, err := getSlice(c.v)
	if err != nil {
		c.t.Fatal(err)
		return c
	}

	if a.Kind() == reflect.String {
		s, err := getString(values...)
		if err != nil {
			c.t.Fatal("argument must be string or array of strings")
			return c
		}
		c.expect(strings.HasPrefix(a.String(), s), "expected '%v' to start with '%v'", c.v, s)
		return c
	}

	if len(values) > a.Len() {
		c.expect(false, "expected %v to start with %v", c.v, values)
		return c
	}

	for i := 0; i < len(values); i++ {
		if !reflect.DeepEqual(a.Index(i).Interface(), values[i]) {
			c.expect(false, "expected %v to start with %v", c.v, values)
			return c
		}
	}

	return c
}

// EndWith checks that subject end with provided elements
func (c *Challenge) EndWith(values ...interface{}) *Challenge {
	c.t.Helper()

	a, err := getSlice(c.v)
	if err != nil {
		c.t.Fatal(err)
		return c
	}

	if a.Kind() == reflect.String {
		s, err := getString(values...)
		if err != nil {
			c.t.Fatal("argument must be string or array of strings")
			return c
		}
		c.expect(strings.HasSuffix(a.String(), s), "expected '%v' to end with '%v'", c.v, s)
		return c
	}

	if len(values) > a.Len() {
		c.expect(false, "expected %v to end with %v", c.v, values)
		return c
	}

	idx := a.Len()-1
	for i := len(values)-1; i >= 0; i-- {
		if !reflect.DeepEqual(a.Index(idx).Interface(), values[i]) {
			c.expect(false, "expected %v to end with %v", c.v, values)
			return c
		}
		idx--
	}

	return c
}

// ContainExactly checks that subject contain exactly all items, regardless of order
func (c *Challenge) ContainExactly(elements ...interface{}) *Challenge {
	c.t.Helper()

	a, err := getSlice(c.v)
	if err != nil {
		c.t.Fatal(err)
		return c
	}

	if a.Kind() == reflect.String {
		s, err := getString(elements...)
		if err != nil {
			c.t.Fatal("argument must be string or array of strings")
			return c
		}
		c.expect(a.String() == s, "expected '%v' to equal '%v'", c.v, s)
		return c
	}

	if len(elements) != a.Len() {
		c.expect(false, "expected %v to equal %v", c.v, elements)
		return c
	}

	for i := 0; i < a.Len(); i++ {
		if !reflect.DeepEqual(a.Index(i).Interface(), elements[i]) {
			c.expect(false, "expected %v to equal %v", c.v, elements)
			return c
		}
	}

	return c
}
