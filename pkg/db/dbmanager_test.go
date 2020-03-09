package db

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/db/memory"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
)

func TestProjectAdd(t *testing.T) {
	prjHandler, _ := memory.NewProjectHandler()
	mgr := &Manager{
		project: prjHandler,
	}

	prjInfo := &model.ProjectInfo{
		Name:      "test-project",
		CreatedAt: time.Now(),
		TokenConfig: &model.TokenConfig{
			AccessTokenLifeSpan:  1,
			RefreshTokenLifeSpan: 1,
			SigningAlgorithm:     "RS256",
		},
	}

	// Test Correct Project
	if err := mgr.ProjectAdd(prjInfo); err != nil {
		t.Errorf("Failed to add correct project: %v", err)
	}

	// Test Duplicate Project Name
	err := mgr.ProjectAdd(prjInfo)
	if errors.Cause(err) != model.ErrProjectAlreadyExists {
		t.Errorf("Expect error is %v, but got %v", model.ErrProjectAlreadyExists, err)
	}

	// TODO(Test Transaction)
}