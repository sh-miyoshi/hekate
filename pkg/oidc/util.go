package oidc

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/oidc/token"
)

// NewAuthRequest ...
func NewAuthRequest(values url.Values) *AuthRequest {
	request := values.Get("request")
	// TODO(parse request and set params)

	maxAge, _ := strconv.Atoi(values.Get("max_age"))
	prompt := []string{}
	if values.Get("prompt") != "" {
		prompt = strings.Split(values.Get("prompt"), " ")
	}
	responseTypes := strings.Split(values.Get("response_type"), " ")

	resMode := values.Get("response_mode")
	if resMode == "" {
		resMode = "fragment"
		if len(responseTypes) == 1 {
			if responseTypes[0] == "code" || responseTypes[0] == "none" {
				resMode = "query"
			}
		}
	}

	return &AuthRequest{
		Scope:               values.Get("scope"),
		ResponseType:        responseTypes,
		ClientID:            values.Get("client_id"),
		RedirectURI:         values.Get("redirect_uri"),
		State:               values.Get("state"),
		Nonce:               values.Get("nonce"),
		Prompt:              prompt,
		MaxAge:              int64(maxAge),
		ResponseMode:        resMode,
		CodeChallenge:       values.Get("code_challenge"),
		CodeChallengeMethod: values.Get("code_challenge_method"),
		Request:             request,
		IDTokenHint:         values.Get("id_token_hint"),
	}
}

// CreateLoggedInResponse ...
func CreateLoggedInResponse(session *model.LoginSession, state, tokenIssuer string) (*http.Request, *errors.Error) {
	values := url.Values{}
	if state != "" {
		values.Set("state", state)
	}

	for _, typ := range session.ResponseType {
		switch typ {
		case "code":
			code := uuid.New().String()
			session.Code = code
			values.Set("code", code)
		case "id_token":
			prj, err := db.GetInst().ProjectGet(session.ProjectName)
			if err != nil {
				return nil, errors.Append(err, "Failed to get token lifespan in project")
			}

			audiences := []string{session.UserID, session.ClientID}
			tokenReq := token.Request{
				Issuer:          tokenIssuer,
				ExpiresIn:       int64(prj.TokenConfig.AccessTokenLifeSpan),
				ProjectName:     session.ProjectName,
				UserID:          session.UserID,
				Nonce:           session.Nonce,
				EndUserAuthTime: session.LoginDate,
			}
			tkn, err := token.GenerateIDToken(audiences, tokenReq)
			if err != nil {
				return nil, errors.Append(err, "Failed to generate id token")
			}
			values.Set("id_token", tkn)
		case "token":
			prj, err := db.GetInst().ProjectGet(session.ProjectName)
			if err != nil {
				return nil, errors.Append(err, "Failed to get token lifespan in project")
			}

			audiences := []string{session.UserID, session.ClientID}
			tokenReq := token.Request{
				Issuer:      tokenIssuer,
				ExpiresIn:   int64(prj.TokenConfig.AccessTokenLifeSpan),
				ProjectName: session.ProjectName,
				UserID:      session.UserID,
			}
			tkn, err := token.GenerateAccessToken(audiences, tokenReq)
			if err != nil {
				return nil, errors.Append(err, "Failed to generate access token")
			}
			values.Set("access_token", tkn)
		default:
			return nil, errors.New("Unknown response type", "Unknown response type %s", typ)
		}
	}

	req, err := http.NewRequest("GET", session.RedirectURI, nil)
	if err != nil {
		return nil, errors.New("Internal server error", "Failed to create response: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if session.ResponseMode == "query" {
		req.URL.RawQuery = values.Encode()
	} else if session.ResponseMode == "fragment" {
		req.URL.Fragment = values.Encode()
	} else {
		return nil, errors.New("Internal server error", "Invalid response mode %s is specified", session.ResponseMode)
	}

	return req, nil
}
