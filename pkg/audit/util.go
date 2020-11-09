package audit

import (
	"net/url"
	"strconv"
	"time"

	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// ParseQuery ...
func ParseQuery(queries *url.Values) (*Query, *errors.Error) {
	now := time.Now()
	res := &Query{
		FromDate: now.AddDate(0, 0, -1),
		ToDate:   now,
		Offset:   0,
	}

	var err error
	if queries.Get("from_date") != "" {
		res.FromDate, err = time.Parse(time.RFC3339, queries.Get("from_date"))
		if err != nil {
			return nil, errors.New("Failed to parse from_date", "Failed to parse from_date: %v", err)
		}
	}
	if queries.Get("to_date") != "" {
		res.ToDate, err = time.Parse(time.RFC3339, queries.Get("to_date"))
		if err != nil {
			return nil, errors.New("Failed to parse to_date", "Failed to parse to_date: %v", err)
		}
	}
	if queries.Get("offset") != "" {
		ofs, err := strconv.Atoi(queries.Get("offset"))
		if err != nil {
			return nil, errors.New("Failed to parse offset", "Failed to parse offset: %v", err)
		}
		if ofs < 0 {
			return nil, errors.New("Offset must be a non-negative", "Offset must be a non-negative, but got %d", ofs)
		}
		res.Offset = uint(ofs)
	}

	return res, nil
}
