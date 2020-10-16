package errors

import (
	"testing"
)

func TestContains(t *testing.T) {
	err1 := New("1", "1")
	err2 := Append(err1, "2")
	err3 := New("3", "3")

	if !Contains(err1, err1) {
		t.Errorf("Unexpect result: Err1 does not contain Err1")
	}

	if !Contains(err2, err1) {
		t.Errorf("Unexpect result: Err2 does not contain Err1")
	}

	if Contains(err2, err3) {
		t.Errorf("Unexpect result: Err2 contains Err3")
	}

	if Contains(nil, nil) {
		t.Errorf("Unexpect result: nil contains nil")
	}

	if Contains(err1, nil) {
		t.Errorf("Unexpect result: Err1 contains nil")
	}

	if Contains(nil, err1) {
		t.Errorf("Unexpect result: nil contains Err1")
	}
}
