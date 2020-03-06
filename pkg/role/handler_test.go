package role

import (
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
