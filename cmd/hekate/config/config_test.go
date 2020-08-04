package config

import (
	"os"
	"testing"
)

func TestSetEnvVar(t *testing.T) {
	k1 := "value1"
	k2 := "value2"
	os.Setenv("KEY1", "new_value")
	setEnvVar("KEY1", &k1)
	setEnvVar("KEY2", &k2)

	if k1 != "new_value" {
		t.Errorf("setEnvVar failed. expect new_value, but got %s", k1)
	}
	if k2 != "value2" {
		t.Errorf("setEnvVar failed. expect value2, but got %s", k2)
	}
}

func TestGetConfigFileName(t *testing.T) {
	tt := []struct {
		args     []string
		expectOK bool
		fname    string
	}{
		{
			args:     []string{},
			expectOK: true,
			fname:    "",
		},
		{
			args:     []string{"--config=test.yaml"},
			expectOK: true,
			fname:    "test.yaml",
		},
		{
			args:     []string{"-config=test.yaml"},
			expectOK: true,
			fname:    "test.yaml",
		},
		{
			args:     []string{"--config", "test.yaml"},
			expectOK: true,
			fname:    "test.yaml",
		},
		{
			args:     []string{"-config", "test.yaml"},
			expectOK: true,
			fname:    "test.yaml",
		},
		{
			args:     []string{"--test=test.yaml"},
			expectOK: true,
			fname:    "",
		},
		{
			args:     []string{"--config"},
			expectOK: false,
		},
		{
			args:     []string{"--config", "--test"},
			expectOK: false,
		},
	}

	for _, tc := range tt {
		f, err := getConfigFileName(tc.args)
		if tc.expectOK && err != nil {
			t.Errorf("getConfigFileName returns wrong response. input: %v, got %v, want nil", tc.args, err)
		}
		if !tc.expectOK && err == nil {
			t.Errorf("getConfigFileName returns wrong response. input: %v, got nil, but want not nil", tc.args)
		}
		if tc.expectOK && f != tc.fname {
			t.Errorf("getConfigFileName returns wrong file name. expect %s, but got %s", tc.fname, f)
		}
	}
}
