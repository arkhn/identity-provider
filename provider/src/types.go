package main

type hydraConfig struct {
	LoginRequestRoute   string
	ConsentRequestRoute string
}

type BodyAcceptOAuth2Login struct {
	Subject     string `json:"subject"`
	Remember    bool   `json:"remember"`
	RememberFor int    `json:"remember_for"`
	Acr         string `json:"acr"`
}

type BodyAcceptOAuth2Consent struct {
	GrantScope               []string `json:"grant_scope"`
	GrantAccessTokenAudience []string `json:"grant_access_token_audience"`
	Remember                 bool     `json:"remember"`
	RememberFor              int      `json:"remember_for"`
	Session                  struct{} `json:"session"` // TODO
}

type RedirectResp struct {
	RedirectTo string `json:"redirect_to"`
}
