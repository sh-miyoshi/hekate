package auditapi

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/hekate/pkg/audit"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	jwthttp "github.com/sh-miyoshi/hekate/pkg/http"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/role"
)

// AuditGetHandler ...
//   require role: read-project
func AuditGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeRead); err != nil {
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Get data form Query
	queries := r.URL.Query()
	logger.Debug("Query: %v", queries)

	now := time.Now()
	fromDate := now.AddDate(0, 0, -1)
	toDate := now
	if queries.Get("from_date") != "" {
		var err error
		fromDate, err = time.Parse(time.RFC3339, queries.Get("from_date"))
		if err != nil {
			logger.Info("Failed to parse from_date: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}
	if queries.Get("to_date") != "" {
		var err error
		toDate, err = time.Parse(time.RFC3339, queries.Get("from_date"))
		if err != nil {
			logger.Info("Failed to parse to_date: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}

	audits, err := audit.GetInst().Get(projectName, fromDate, toDate)
	if err != nil {
		errors.Print(errors.Append(err, "Failed to get audit events"))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res := []*AuditGetResponse{}
	for _, a := range audits {
		res = append(res, &AuditGetResponse{
			Time:         a.Time,
			ResourceType: a.ResourceType,
			Method:       a.Method,
			Path:         a.Path,
			IsSuccess:    a.IsSuccess,
			Message:      a.Message,
		})
	}

	jwthttp.ResponseWrite(w, "AuditGetHandler", res)
}
