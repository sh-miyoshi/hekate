package auditapi

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/hekate/pkg/audit"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	jwthttp "github.com/sh-miyoshi/hekate/pkg/http"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/role"
)

type auditRequest struct {
	FromDate time.Time
	ToDate   time.Time
	Offset   uint
}

func parseQuery(queries *url.Values) (*auditRequest, *errors.Error) {
	now := time.Now()
	res := &auditRequest{
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

// AuditGetHandler ...
//   require role: read-project
func AuditGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeRead); err != nil {
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		errors.WriteHTTPError(w, "Forbidden", err, http.StatusForbidden)
		return
	}

	// Get data form Query
	queries := r.URL.Query()
	logger.Debug("Query: %v", queries)
	req, err := parseQuery(&queries)
	if err != nil {
		errors.PrintAsInfo(errors.Append(err, "Failed to parse query"))
		errors.WriteHTTPError(w, "Bad Request", err, http.StatusBadRequest)
		return
	}

	audits, err := audit.GetInst().Get(projectName, req.FromDate, req.ToDate, req.Offset)
	if err != nil {
		errors.Print(errors.Append(err, "Failed to get audit events"))
		errors.WriteHTTPError(w, "Internal Server Error", err, http.StatusInternalServerError)
		return
	}

	res := []*AuditGetResponse{}
	for _, a := range audits {
		res = append(res, &AuditGetResponse{
			Time:         a.Time.Format(time.RFC3339),
			ResourceType: a.ResourceType,
			Method:       a.Method,
			Path:         a.Path,
			IsSuccess:    a.IsSuccess,
			Message:      a.Message,
		})
	}

	jwthttp.ResponseWrite(w, "AuditGetHandler", res)
}
