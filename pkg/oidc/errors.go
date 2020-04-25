package oidc

import (
	"net/http"
)

var (
	//-------------------------------------
	// Define in RFC6749
	//-------------------------------------

	// ErrUnknownRequest ...
	ErrUnknownRequest = &Error{
		Name: "unknown_error",
		Code: http.StatusBadRequest,
	}

	// ErrRequestForbidden ...
	ErrRequestForbidden = &Error{
		Name: "request_forbidden",
		Code: http.StatusForbidden,
	}

	// ErrInvalidRequest ...
	ErrInvalidRequest = &Error{
		Name: "invalid_request",
		Code: http.StatusBadRequest,
	}

	// ErrUnauthorizedClient ...
	ErrUnauthorizedClient = &Error{
		Name: "unauthorized_client",
		Code: http.StatusBadRequest,
	}

	// ErrAccessDenied ...
	ErrAccessDenied = &Error{
		Name: "access_denied",
		Code: http.StatusForbidden,
	}

	// ErrUnsupportedResponseType ...
	ErrUnsupportedResponseType = &Error{
		Name:        "unsupported_response_type",
		Description: "Given type is unsupported",
		Code:        http.StatusBadRequest,
	}

	// ErrInvalidScope ...
	ErrInvalidScope = &Error{
		Name: "invalid_scope",
		Code: http.StatusBadRequest,
	}

	// ErrServerError ...
	ErrServerError = &Error{
		Name: "server_error",
		Code: http.StatusInternalServerError,
	}

	// ErrTemporarilyUnavailable ...
	ErrTemporarilyUnavailable = &Error{
		Name: "temporarily_unavailable",
		Code: http.StatusServiceUnavailable,
	}

	// ErrUnsupportedGrantType ...
	ErrUnsupportedGrantType = &Error{
		Name: "unsupported_grant_type",
		Code: http.StatusBadRequest,
	}

	// ErrInvalidGrant ...
	ErrInvalidGrant = &Error{
		Name: "invalid_grant",
		Code: http.StatusBadRequest,
	}

	// ErrInvalidClient ...
	ErrInvalidClient = &Error{
		Name:        "invalid_client",
		Description: "Client authentication failed.",
		Code:        http.StatusUnauthorized,
	}

	// ErrInvalidState ...
	ErrInvalidState = &Error{
		Name: "invalid_state",
		Code: http.StatusBadRequest,
	}

	// ErrMisconfiguration ...
	ErrMisconfiguration = &Error{
		Name: "misconfiguration",
		Code: http.StatusInternalServerError,
	}

	// ErrInsufficientEntropy ...
	ErrInsufficientEntropy = &Error{
		Name: "insufficient_entropy",
		Code: http.StatusBadRequest,
	}

	// ErrNotFound ...
	ErrNotFound = &Error{
		Name: "not_found",
		Code: http.StatusNotFound,
	}

	// ErrRequestUnauthorized ...
	ErrRequestUnauthorized = &Error{
		Name: "request_unauthorized",
		Code: http.StatusUnauthorized,
	}

	// ErrTokenSignatureMismatch ...
	ErrTokenSignatureMismatch = &Error{
		Name: "token_signature_mismatch",
		Code: http.StatusBadRequest,
	}

	// ErrInvalidTokenFormat ...
	ErrInvalidTokenFormat = &Error{
		Name: "invalid_token_format",
		Code: http.StatusBadRequest,
	}

	// ErrTokenExpired ...
	ErrTokenExpired = &Error{
		Name: "token_expired",
		Code: http.StatusUnauthorized,
	}

	// ErrScopeNotGranted ...
	ErrScopeNotGranted = &Error{
		Name: "scope_not_granted",
		Code: http.StatusForbidden,
	}

	// ErrTokenClaim ...
	ErrTokenClaim = &Error{
		Name: "token_claim",
		Code: http.StatusUnauthorized,
	}

	// ErrInactiveToken ...
	ErrInactiveToken = &Error{
		Name: "token_inactive",
		Code: http.StatusUnauthorized,
	}

	// ErrRevokationClientMismatch ...
	ErrRevokationClientMismatch = &Error{
		Name: "revokation_client_mismatch",
		Code: http.StatusBadRequest,
	}

	// ErrLoginRequired ...
	ErrLoginRequired = &Error{
		Name: "login_required",
		Code: http.StatusBadRequest,
	}

	// ErrInteractionRequired ...
	ErrInteractionRequired = &Error{
		Name: "interaction_required",
		Code: http.StatusBadRequest,
	}

	// ErrConsentRequired ...
	ErrConsentRequired = &Error{
		Name: "consent_required",
		Code: http.StatusBadRequest,
	}

	// ErrRequestNotSupported ...
	ErrRequestNotSupported = &Error{
		Name: "request_not_supported",
		Code: http.StatusBadRequest,
	}

	// ErrRequestURINotSupported ...
	ErrRequestURINotSupported = &Error{
		Name: "request_uri_not_supported",
		Code: http.StatusBadRequest,
	}

	// ErrRegistrationNotSupported ...
	ErrRegistrationNotSupported = &Error{
		Name: "registration_not_supported",
		Code: http.StatusBadRequest,
	}

	// ErrInvalidRequestURI ...
	ErrInvalidRequestURI = &Error{
		Name: "invalid_request_uri",
		Code: http.StatusBadRequest,
	}

	// ErrInvalidRequestObject ...
	ErrInvalidRequestObject = &Error{
		Name: "invalid_request_object",
		Code: http.StatusBadRequest,
	}

	//-------------------------------------
	// Define in OAuth 2.0 Token Revocation
	//-------------------------------------

	// ErrUnsupportedTokenType ...
	ErrUnsupportedTokenType = &Error{
		Name: "unsupported_token_type",
		Code: http.StatusBadRequest,
	}

	//-------------------------------------
	// Define in Open ID Connect
	//-------------------------------------

	// ErrAccountSelectionRequired ...
	ErrAccountSelectionRequired = &Error{
		Name: "account_selection_required",
		Code: http.StatusBadRequest,
	}
)
