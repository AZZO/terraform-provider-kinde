package kinde

type ApiRequestApplicationUpdate struct {
	Name         string   `json:"name"`
	LogoutUris   []string `json:"logout_uris"`
	RedirectUris []string `json:"redirect_uris"`
}
