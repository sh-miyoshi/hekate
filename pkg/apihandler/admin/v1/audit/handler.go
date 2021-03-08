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
		errors.WriteToHTTP(w, errors.ErrUnpermitted, 0, "")
		return
	}

	// Get data form Query
	queries := r.URL.Query()
	logger.Debug("Query: %v", queries)
	req, err := audit.ParseQuery(&queries)
	if err != nil {
		errors.PrintAsInfo(errors.Append(err, "Failed to parse query"))
		errors.WriteToHTTP(w, err, http.StatusBadRequest, "")
		return
	}

	audits, err := audit.GetInst().Get(projectName, *req)
	if err != nil {
		errors.Print(errors.Append(err, "Failed to get audit events"))
		errors.WriteToHTTP(w, err, http.StatusInternalServerError, "")
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
