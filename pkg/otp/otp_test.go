package otp

import (
	"encoding/base32"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
)

func TestVerify(t *testing.T) {
	privateKey := []byte("12345678901234567890")

	user := &model.UserInfo{
		ID:          uuid.New().String(),
		ProjectName: "master",
		Name:        "admin",
		OTPInfo: model.OTPInfo{
			ID:         uuid.New().String(),
			PrivateKey: base32.StdEncoding.EncodeToString(privateKey),
			Enabled:    true,
		},
	}

	tt := []struct {
		TimeSec int
		Expect  string
	}{
		{TimeSec: 59, Expect: "287082"},
		{TimeSec: 1111111109, Expect: "081804"},
		{TimeSec: 1111111111, Expect: "050471"},
		{TimeSec: 1234567890, Expect: "005924"},
		{TimeSec: 2000000000, Expect: "279037"},
		{TimeSec: 20000000000, Expect: "353130"},
	}

	for _, tc := range tt {
		if err := Verify(time.Unix(int64(tc.TimeSec), 0), user, tc.Expect); err != nil {
			t.Errorf("Failed to verify user code: %v", err)
		}
	}
}

func TestTruncate(t *testing.T) {
	input := []byte{0x1f, 0x86, 0x98, 0x69, 0x0e, 0x02, 0xca, 0x16, 0x61, 0x85, 0x50, 0xef, 0x7f, 0x19, 0xda, 0x8e, 0x94, 0x5b, 0x55, 0x5a}
	expect := 872921

	res := truncate(input)
	if res != expect {
		t.Errorf("truncate method return %d, but want %d", res, expect)
	}
}
