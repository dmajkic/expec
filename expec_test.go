package expec

import "testing"

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