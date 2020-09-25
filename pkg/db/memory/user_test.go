package memory

import (
	"testing"

	"github.com/sh-miyoshi/hekate/pkg/db/model"
)

func TestFilterUserList(t *testing.T) {
	tt := []struct {
		UserName   string
		FilterName string
		ExpectNum  int
	}{
		{
			UserName:   "admin",
			FilterName: "admin",
			ExpectNum:  1,
		},
		{
			UserName:   "admin",
			FilterName: "",
			ExpectNum:  1,
		},
		{
			UserName:   "admin",
			FilterName: "fakeadmin",
			ExpectNum:  0,
		},
		{
			UserName:   "admin",
			FilterName: "adminfake",
			ExpectNum:  0,
		},
		// TODO add test for filtering by ID
	}

	for _, tc := range tt {
		const project = "master"

		data := []*model.UserInfo{
			{
				ProjectName: project,
				Name:        tc.UserName,
			},
		}

		filter := &model.UserFilter{
			Name: tc.FilterName,
		}

		res := filterUserList(data, project, filter)
		if len(res) != tc.ExpectNum {
			t.Errorf("Filter User List Failed: expect num: %d, but got %d, %v", tc.ExpectNum, len(res), res)
		}
	}
}
