package errors

import (
	"net/http"
)

var (
	//-------------------------------------
	// Define in RFC6749
	//-------------------------------------

	// ErrUnknownRequest ...
	ErrUnknownRequest = &Error{
		publicMsg:        "unknown_error",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrRequestForbidden ...
	ErrRequestForbidden = &Error{
		publicMsg:        "request_forbidden",
		httpResponseCode: http.StatusForbidden,
	}

	// ErrInvalidRequest ...
	ErrInvalidRequest = &Error{
		publicMsg:        "invalid_request",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrUnauthorizedClient ...
	ErrUnauthorizedClient = &Error{
		publicMsg:        "unauthorized_client",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrAccessDenied ...
	ErrAccessDenied = &Error{
		publicMsg:        "access_denied",
		httpResponseCode: http.StatusForbidden,
	}

	// ErrUnsupportedResponseType ...
	ErrUnsupportedResponseType = &Error{
		publicMsg:        "unsupported_response_type",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrInvalidScope ...
	ErrInvalidScope = &Error{
		publicMsg:        "invalid_scope",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrServerError ...
	ErrServerError = &Error{
		publicMsg:        "server_error",
		httpResponseCode: http.StatusInternalServerError,
	}

	// ErrTemporarilyUnavailable ...
	ErrTemporarilyUnavailable = &Error{
		publicMsg:        "temporarily_unavailable",
		httpResponseCode: http.StatusServiceUnavailable,
	}

	// ErrUnsupportedGrantType ...
	ErrUnsupportedGrantType = &Error{
		publicMsg:        "unsupported_grant_type",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrInvalidGrant ...
	ErrInvalidGrant = &Error{
		publicMsg:        "invalid_grant",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrInvalidClient ...
	ErrInvalidClient = &Error{
		publicMsg:        "invalid_client",
		httpResponseCode: http.StatusUnauthorized,
	}

	// ErrInvalidState ...
	ErrInvalidState = &Error{
		publicMsg:        "invalid_state",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrMisconfiguration ...
	ErrMisconfiguration = &Error{
		publicMsg:        "misconfiguration",
		httpResponseCode: http.StatusInternalServerError,
	}

	// ErrInsufficientEntropy ...
	ErrInsufficientEntropy = &Error{
		publicMsg:        "insufficient_entropy",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrNotFound ...
	ErrNotFound = &Error{
		publicMsg:        "not_found",
		httpResponseCode: http.StatusNotFound,
	}

	// ErrRequestUnauthorized ...
	ErrRequestUnauthorized = &Error{
		publicMsg:        "request_unauthorized",
		httpResponseCode: http.StatusUnauthorized,
	}

	// ErrTokenSignatureMismatch ...
	ErrTokenSignatureMismatch = &Error{
		publicMsg:        "token_signature_mismatch",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrInvalidTokenFormat ...
	ErrInvalidTokenFormat = &Error{
		publicMsg:        "invalid_token_format",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrTokenExpired ...
	ErrTokenExpired = &Error{
		publicMsg:        "token_expired",
		httpResponseCode: http.StatusUnauthorized,
	}

	// ErrScopeNotGranted ...
	ErrScopeNotGranted = &Error{
		publicMsg:        "scope_not_granted",
		httpResponseCode: http.StatusForbidden,
	}

	// ErrTokenClaim ...
	ErrTokenClaim = &Error{
		publicMsg:        "token_claim",
		httpResponseCode: http.StatusUnauthorized,
	}

	// ErrInactiveToken ...
	ErrInactiveToken = &Error{
		publicMsg:        "token_inactive",
		httpResponseCode: http.StatusUnauthorized,
	}

	// ErrRevokationClientMismatch ...
	ErrRevokationClientMismatch = &Error{
		publicMsg:        "revokation_client_mismatch",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrLoginRequired ...
	ErrLoginRequired = &Error{
		publicMsg:        "login_required",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrInteractionRequired ...
	ErrInteractionRequired = &Error{
		publicMsg:        "interaction_required",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrConsentRequired ...
	ErrConsentRequired = &Error{
		publicMsg:        "consent_required",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrRequestNotSupported ...
	ErrRequestNotSupported = &Error{
		publicMsg:        "request_not_supported",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrRequestURINotSupported ...
	ErrRequestURINotSupported = &Error{
		publicMsg:        "request_uri_not_supported",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrRegistrationNotSupported ...
	ErrRegistrationNotSupported = &Error{
		publicMsg:        "registration_not_supported",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrInvalidRequestURI ...
	ErrInvalidRequestURI = &Error{
		publicMsg:        "invalid_request_uri",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrInvalidRequestObject ...
	ErrInvalidRequestObject = &Error{
		publicMsg:        "invalid_request_object",
		httpResponseCode: http.StatusBadRequest,
	}

	//-------------------------------------
	// Define in OAuth 2.0 Token Revocation
	//-------------------------------------

	// ErrUnsupportedTokenType ...
	ErrUnsupportedTokenType = &Error{
		publicMsg:        "unsupported_token_type",
		httpResponseCode: http.StatusBadRequest,
	}

	//-------------------------------------
	// Define in Open ID Connect
	//-------------------------------------

	// ErrAccountSelectionRequired ...
	ErrAccountSelectionRequired = &Error{
		publicMsg:        "account_selection_required",
		httpResponseCode: http.StatusBadRequest,
	}

	//-------------------------------------
	// Original
	//-------------------------------------

	// ErrSessionExpired ...
	ErrSessionExpired = &Error{
		publicMsg:        "already_session_expired",
		httpResponseCode: http.StatusBadRequest,
	}
)
