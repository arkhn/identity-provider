package main

type hydraConfig struct {
	LoginRequestRoute   string
	ConsentRequestRoute string
}

type BodyAcceptOAuth2Login struct {
	Acr         string `json:"acr"`
	Remember    bool   `json:"remember"`
	RememberFor int    `json:"remember_for"`
	Subject     string `json:"subject"`
}

type IdTokenClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type SessionInfo struct {
	// AccessToken string        `json:"access_token"` // TODO?
	IdToken IdTokenClaims `json:"id_token"`
}

type BodyAcceptOAuth2Consent struct {
	GrantScope               []string    `json:"grant_scope"`
	GrantAccessTokenAudience []string    `json:"grant_access_token_audience"`
	Remember                 bool        `json:"remember"`
	RememberFor              int         `json:"remember_for"`
	Session                  SessionInfo `json:"session"` // TODO
}

type RedirectResp struct {
	RedirectTo string `json:"redirect_to"`
}
