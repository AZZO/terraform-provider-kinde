package kinde_client

type ApplicationResource struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type requestApplicationUpdate struct {
	Name         string   `json:"name"`
	LogoutUris   []string `json:"logout_uris"`
	RedirectUris []string `json:"redirect_uris"`
}

type responseApplicationGet struct {
	Code        string              `json:"code"`
	Message     string              `json:"message"`
	Application ApplicationResource `json:"application"`
}

type requestApplicationCreate struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type responseApplicationCreate struct {
	Code        string              `json:"code"`
	Message     string              `json:"message"`
	Application ApplicationResource `json:"application"`
}
