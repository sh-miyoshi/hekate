package login

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/util"
)

func TestVerify(t *testing.T) {
	const userName = "test-user"
	const password = "test-password"

	// Initialize test DB
	db.InitDBManager("memory", "")
	db.GetInst().ProjectAdd(&model.ProjectInfo{
		Name: "master",
		TokenConfig: &model.TokenConfig{
			AccessTokenLifeSpan:  model.DefaultAccessTokenExpiresTimeSec,
			RefreshTokenLifeSpan: model.DefaultRefreshTokenExpiresTimeSec,
			SigningAlgorithm:     "RS256",
		},
	})
	db.GetInst().UserAdd("master", &model.UserInfo{
		ID:           uuid.New().String(),
		ProjectName:  "master",
		Name:         userName,
		PasswordHash: util.CreateHash(password),
	})

	tt := []struct {
		projectName string
		userName    string
		password    string
		expectOk    bool
	}{
		{
			"master",
			userName,
			password,
			true,
		},
		{
			"wrong-project",
			userName,
			password,
			false,
		},
		{
			"master",
			"wrong-user",
			password,
			false,
		},
		{
			"master",
			userName,
			"wrong-password",
			false,
		},
	}

	for _, tc := range tt {
		_, err := UserVerifyByPassword(tc.projectName, tc.userName, tc.password)
		if tc.expectOk && err != nil {
			t.Errorf("Verify returns wrong response. got %v, want nil", err)
		}
		if !tc.expectOk && err == nil {
			t.Errorf("Verify returns wrong response. got nil, but want not nil")
		}
	}
}
