package errors

import (
	"net/http"
)

var (
	//-------------------------------------
	// Define in RFC6749
	//-------------------------------------

	// ErrInvalidRequest ...
	ErrInvalidRequest = &Error{
		publicMsg:        "invalid_request",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrInvalidClient ...
	ErrInvalidClient = &Error{
		publicMsg:        "invalid_client",
		httpResponseCode: http.StatusUnauthorized,
	}

	// ErrInvalidGrant ...
	ErrInvalidGrant = &Error{
		publicMsg:        "invalid_grant",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrUnauthorizedClient ...
	ErrUnauthorizedClient = &Error{
		publicMsg:        "unauthorized_client",
		httpResponseCode: http.StatusBadRequest,
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

	//-------------------------------------
	// Define in RFC6750
	//-------------------------------------

	// ErrInvalidToken ...
	ErrInvalidToken = &Error{
		publicMsg:        "invalid_token",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrInsufficientScope ...
	ErrInsufficientScope = &Error{
		publicMsg:        "insufficient_scope",
		httpResponseCode: http.StatusForbidden,
	}

	//-------------------------------------
	// Define in Open ID Connect
	//-------------------------------------

	// ErrAccountSelectionRequired ...
	ErrAccountSelectionRequired = &Error{
		publicMsg:        "account_selection_required",
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
	// RFC 8628
	//-------------------------------------

	// ErrAuthorizationPending ...
	ErrAuthorizationPending = &Error{
		publicMsg:        "authorization_pending",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrSlowDown ...
	ErrSlowDown = &Error{
		publicMsg:        "slow_down",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrAccessDenied ...
	ErrAccessDenied = &Error{
		publicMsg:        "access_denied",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrExpiredToken ...
	ErrExpiredToken = &Error{
		publicMsg:        "expired_token",
		httpResponseCode: http.StatusBadRequest,
	}

	//-------------------------------------
	// Original
	//-------------------------------------

	// ErrRequestUnauthorized ...
	ErrRequestUnauthorized = &Error{
		publicMsg:        "request_unauthorized",
		httpResponseCode: http.StatusUnauthorized,
	}

	// ErrSessionExpired ...
	ErrSessionExpired = &Error{
		publicMsg:        "already_session_expired",
		httpResponseCode: http.StatusBadRequest,
	}

	// ErrProjectNotFound ...
	ErrProjectNotFound = &Error{
		publicMsg:        "project_not_found",
		httpResponseCode: http.StatusNotFound,
	}

	// ErrUnpermitted ...
	ErrUnpermitted = &Error{
		publicMsg:        "no_permission",
		httpResponseCode: http.StatusForbidden,
	}
)
