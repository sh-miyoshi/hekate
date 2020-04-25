package role

import (
	"fmt"
	"testing"
)

func TestAuthorize(t *testing.T) {
	roles := []string{
		"read-test1",
		"write-test1",
		"read-test2",
	}

	tt := []struct {
		targetResource Resource
		roleType       Type
		expectSuccess  bool
	}{
		{
			Resource{"test1"},
			Type{"read"},
			true,
		},
		{
			Resource{"test2"},
			Type{"read"},
			true,
		},
		{
			Resource{"test1"},
			Type{"write"},
			true,
		},
		{
			Resource{"test2"},
			Type{"write"},
			false,
		},
		{
			Resource{"test"},
			Type{"read"},
			false,
		},
	}

	for _, target := range tt {
		res := target.targetResource
		typ := target.roleType
		result := Authorize(roles, res, typ)
		if result != target.expectSuccess {
			t.Errorf("Authorize role %s-%s returns wrong status. got %v, want %v", res, typ, result, target.expectSuccess)
		}
	}
}

func TestParse(t *testing.T) {
	// Initialize DB
	InitHandler()

	tt := []struct {
		input    string
		expectOk bool
	}{
		{
			"read-cluster",
			true,
		},
		{
			"badtype-cluster",
			false,
		},
		{
			"read-badresource",
			false,
		},
		{
			"readcluster",
			false,
		},
	}

	for _, tc := range tt {
		res, typ, ok := GetInst().Parse(tc.input)
		if tc.expectOk != ok {
			t.Errorf("Parse %s returns wrong response. got %v, want %v", tc.input, ok, tc.expectOk)
		}
		if ok {
			v := fmt.Sprintf("%s-%s", typ.String(), res.String())
			if v != tc.input {
				t.Errorf("Parse %s returns wrong response. got %s, want %s", tc.input, v, tc.input)
			}
		}
	}
}
