package oidc

type authInfo struct {
	Scope        string // scope(REQUIRED)
	ResponseType string // response_type(REQUIRED)
	ClientID     string // client_id(REQUIRED)
	RedirectURI  string // redirect_uri(REQUIRED)
	State        string // state(RECOMMENDED)

	// TODO(implement this)
	// ResponseMode string // response_mode(OPTIONAL)
	// Nonce string // nonce(OPTIONAL)
	// Display string // display(OPTIONAL)
	// Prompt string // prompt(OPTIONAL)
	// MaxAge string // max_age(OPTIONAL)
	// UILocales string // ui_locales(OPTIONAL)
	// IDTokenHint string // id_token_hint(OPTIONAL)
	// LoginHint string // login_hint(OPTIONAL)
	// ACRValues string // acr_values(OPTIONAL)
}

// Validate ...
func (ai *authInfo) Validate() error {
	return nil
}
