package util

import (
	"testing"
)

func TestInt2bytes(t *testing.T) {
	tt := []struct {
		input  uint64
		expect []byte
	}{
		{
			0,
			[]byte{0},
		},
		{
			3,
			[]byte{3},
		},
		{
			259,
			[]byte{1, 3},
		},
		{
			876548754787,
			[]byte{204, 22, 96, 141, 99},
		},
	}

	for _, tc := range tt {
		res := Int2bytes(tc.input)
		if len(res) != len(tc.expect) {
			t.Errorf("Int2bytes returns wrong response. got %v, want %v", res, tc.expect)
		}
		for i := 0; i < len(res); i++ {
			if res[i] != tc.expect[i] {
				t.Errorf("Int2bytes returns wrong response. got %v, want %v", res, tc.expect)
				break
			}
		}
	}
}
