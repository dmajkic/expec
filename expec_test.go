package expec

import (
	"errors"
	"os"
	"testing"
)

func TestExpec (t *testing.T) {
	v := "test"
	if Expec(t, v).v != v {
		t.Errorf("expected 'test', got %v", Expec(t, v).v)
	}

	if Expec(t, v).t != t {
		t.Errorf("wrong t")
	}

	cTo := Expec(t, v).To
	if cTo.t != t {
		t.Errorf("wrong t")
	}

	if cTo.not != false {
		t.Errorf("wrong challenge, should not be !")
	}

	cNotTo := Expec(t, v).NotTo
	if cNotTo.t != t {
		t.Errorf("wrong t")
	}

	if cNotTo.not != true {
		t.Errorf("wrong challenge should be !")
	}
}

func TestSubject (t *testing.T) {
	v := "test"

	if Expec(t, v).To.v != "test" {
		t.Errorf("expected 'test', got %v", Expec(t, v).To.v)
	}

	if Expec(t, v).NotTo.v != "test" {
		t.Errorf("expected 'test', got %v", Expec(t, v).NotTo.v)
	}
}

func TestChallenge_Be(t *testing.T) {
	v := "test"
	Expec(t, v).To.Be("test")

	type s struct {a int; b int; c int}
	sv := s{a:1,b:2,c:3}
	Expec(t, s{1,2,3}).To.Be(sv)

	Expec(t, s{1,2,3}).NotTo.Be([]int{1,2,3})

	var n *int
	Expec(t, n).To.BeNil()
	Expec(t, nil).To.BeNil()
}

func TestChallenge_BeA(t *testing.T) {
	var e error
	Expec(t, e).To.BeA(error(nil))
}

func TestChallenge_BeAn(t *testing.T) {
	var e error
	Expec(t, e).To.BeA(error(nil))
}

func TestChallenge_BeFalse(t *testing.T) {
	v := false
	Expec(t, v).To.BeFalse()

	Expec(t, true).NotTo.BeFalse()
	Expec(t, 0).NotTo.BeFalse()
	Expec(t, "").NotTo.BeFalse()
	Expec(t, nil).NotTo.BeFalse()
}

func TestChallenge_BeFalsy(t *testing.T) {
	Expec(t, false).To.BeFalsy()
	Expec(t, nil).To.BeFalsy()

	Expec(t, true).NotTo.BeFalsy()
}

func TestChallenge_BeNil(t *testing.T) {
	Expec(t, nil).To.BeNil()
	Expec(t, error(nil)).To.BeNil()

	var v int
	Expec(t, v).NotTo.BeNil()
	Expec(t, struct{}{}).NotTo.BeNil()

	var p *int
	Expec(t, p).To.BeNil()
}

func TestChallenge_BeTrue(t *testing.T) {
	v := true
	Expec(t, v).To.BeTrue()

	Expec(t, false).NotTo.BeTrue()
	Expec(t, 0).NotTo.BeTrue()
	Expec(t, "").NotTo.BeTrue()
	Expec(t, nil).NotTo.BeTrue()
}

func TestChallenge_BeTruthy(t *testing.T) {
	Expec(t, false).NotTo.BeTruthy()
	Expec(t, nil).NotTo.BeTruthy()

	Expec(t, true).To.BeTruthy()
}

func TestChallenge_Eq(t *testing.T) {
	v := 1
	Expec(t, v).To.Eq(1)

	type s struct {a int; b int; c *int}
	tre := 3
	vtr := 3
	sv := s{a:1,b:2,c:&tre}

	Expec(t, s{1,2,&tre}).To.Eq(sv)
	Expec(t, s{1,2,&vtr}).NotTo.Eq([]int{1,2,3})
}

func TestChallenge_Eql(t *testing.T) {
	v := 1
	Expec(t, v).To.Eql(1)

	type s struct {a int; b int; c *int}
	tre := 3
	vtr := 3
	sv := s{a:1,b:2,c:&tre}

	Expec(t, s{1,2,&tre}).To.Eql(sv)
	Expec(t, s{1,2,&vtr}).NotTo.Eql([]int{1,2,3})
}

func TestChallenge_Implement(t *testing.T) {
	Expec(t, errors.New("error raised")).To.Implement((*error)(nil))

	var e error
	Expec(t, e).NotTo.Implement((*error)(nil))
}

func TestChallenge_Include(t *testing.T) {
	Expec(t, []int{1,2,3,4}).To.Include(2)
	Expec(t, []string{"1","2","3","4"}).To.Include("4")
	Expec(t, []interface{}{"1",2,3.33,"4"}).To.Include(3.33)
}

func TestChallenge_Match(t *testing.T) {
	Expec(t, "Something nice").To.Match("nice$")
	Expec(t, "Something nice").To.Match("^Some")
	Expec(t, "Something nice").NotTo.Match("^some")
	Expec(t, "Something nice").To.Match("(?i)^some")
}

func TestChallenge_RaiseError(t *testing.T) {
	err := errors.New("error raised")
	Expec(t, err).To.RaiseError()
	Expec(t, err).To.RaiseError("error raised")

	err2 := os.ErrClosed
	Expec(t, err2).To.RaiseError(os.ErrClosed)

	Expec(t, err).NotTo.BeNil()
}
