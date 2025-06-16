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

type ApiResource struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Audience string `json:"audience"`
}

type requestApiCreate struct {
	Name     string `json:"name"`
	Audience string `json:"audience"`
}

type responseApiGet struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Api     ApiResource `json:"api"`
}

type responseApiCreate struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Api     ApiResource `json:"api"`
}

type EnvironmentResource struct {
	Code            string `json:"code"`
	Name            string `json:"name"`
	IsDefault       bool   `json:"is_default"`
	IsLive          bool   `json:"is_live"`
	KindeDomain     string `json:"kinde_domain"`
	CustomDomain    string `json:"custom_domain"`
	Logo            string `json:"logo"`
	LogoDark        string `json:"logo_dark"`
	FaviconSvg      string `json:"favicon_svg"`
	FaviconFallback string `json:"favicon_fallback"`
	CreatedOn       string `json:"created_on"`
}

type responseEnvironmentGet struct {
	Code        string              `json:"code"`
	Message     string              `json:"message"`
	Environment EnvironmentResource `json:"environment"`
}
