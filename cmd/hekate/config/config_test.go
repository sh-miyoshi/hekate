package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestInitConfig(t *testing.T) {
	tmpFile, _ := ioutil.TempFile("", "testconf")
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString(`admin_name: admin
admin_password: password
server_port: 8080
server_bind_address: "0.0.0.0"
logfile: ''
debug_mode: true
db:
  type: "memory"
  connection_string: ""
oidc_auth_code_expires_time: 300`)

	// Test loading correct yaml file
	if _, err := InitConfig(tmpFile.Name()); err != nil {
		t.Errorf("Failed to load correct yaml file: %v", err)
	}

	// Test overwrite by env
	admin := "hekate"
	passwd := "securepassword"
	os.Setenv("HEKATE_ADMIN_NAME", admin)
	os.Setenv("HEKATE_ADMIN_PASSWORD", passwd)
	cfg, _ := InitConfig(tmpFile.Name())
	if cfg.AdminName != admin || cfg.AdminPassword != passwd {
		t.Errorf("Failed to overwrite by os.Env: want %s:%s, got %s:%s", admin, passwd, cfg.AdminName, cfg.AdminPassword)
	}
}
