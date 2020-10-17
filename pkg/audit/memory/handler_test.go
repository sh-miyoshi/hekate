package memory

import (
	"testing"
	"time"

	"github.com/sh-miyoshi/hekate/pkg/audit/model"
	"github.com/sh-miyoshi/hekate/pkg/util"
)

func TestGet(t *testing.T) {
	// Set test data
	handler := NewHandler()
	baseTime := time.Now()
	for i := -50; i < 250; i++ {
		tm := baseTime.AddDate(0, 0, i)
		handler.data = append(handler.data, model.Audit{
			ProjectName:  "project1",
			Time:         tm,
			ResourceType: "TEST",
			Method:       "GET",
			Path:         "/",
			IsSuccess:    true,
			Message:      "",
		})
	}

	// Test fromDate and toDate
	fromDate := util.TimeTruncate(baseTime)
	n := 10
	toDate := util.TimeTruncate(baseTime.AddDate(0, 0, n-1))
	res, _ := handler.Get("project1", fromDate, toDate, 0)
	if len(res) != n {
		t.Errorf("Failed to get audit events by filtering date. expect: %d, but got: %d", n, len(res))
	}

	// Test max num
	fromDate = util.TimeTruncate(baseTime.AddDate(0, 0, -50))
	toDate = util.TimeTruncate(baseTime.AddDate(0, 0, 250))
	res, _ = handler.Get("project1", fromDate, toDate, 0)
	if len(res) != model.AuditGetMaxNum {
		t.Errorf("Failed to get audit events by max num. expect: %d, but got: %d", model.AuditGetMaxNum, len(res))
	}

	// Test offset
	fromDate = util.TimeTruncate(baseTime)
	offset := uint(1)
	res, _ = handler.Get("project1", fromDate, toDate, offset)
	if !res[0].Time.Equal(handler.data[150].Time) {
		t.Errorf("Failed to get audit event by offset. expect: %v, but got: %v", handler.data[150].Time, res[0].Time)
	}
	offset = uint(2)
	res, _ = handler.Get("project1", fromDate, toDate, offset)
	if !res[0].Time.Equal(handler.data[250].Time) {
		t.Errorf("Failed to get audit event by offset 2. expect: %v, but got: %v", handler.data[250].Time, res[0].Time)
	}
	if len(res) != 50 {
		t.Errorf("Failed to get audit event num by offset 2. expect: %d, but got: %d", 50, len(res))
	}
	offset = uint(100)
	res, _ = handler.Get("project1", fromDate, toDate, offset)
	if len(res) != 0 {
		t.Errorf("Failed to get audit event by over offset. expect: 0, but got: %d", len(res))
	}
}
