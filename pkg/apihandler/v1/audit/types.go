package auditapi

// AuditGetResponse ...
type AuditGetResponse struct {
	Time         string `json:"time"`
	ResourceType string `json:"resource_type"`
	Method       string `json:"method"`
	Path         string `json:"path"`
	IsSuccess    bool   `json:"success"`
	Message      string `json:"message"`
}
