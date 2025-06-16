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

type responseApplicationDelete struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ApiResource struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Audience string `json:"audience"`
}

type ApiScopeResource struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type requestApiCreate struct {
	Name     string `json:"name"`
	Audience string `json:"audience"`
}

type requestApiScopeCreate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type responseApiGet struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Api     ApiResource `json:"api"`
}

type responseApisGet struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Apis    []ApiResource `json:"apis"`
}

type responseApiCreate struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Api     ApiResource `json:"api"`
}

type responseApiScopesGet struct {
	Code    string             `json:"code"`
	Message string             `json:"message"`
	Scopes  []ApiScopeResource `json:"scopes"`
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

type RoleResource struct {
	Id            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Key           string   `json:"key"`
	IsDefaultRole bool     `json:"is_default_role"`
	Permissions   []string `json:"permissions"`
}

type requestRoleCreate struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Key           string   `json:"key"`
	IsDefaultRole bool     `json:"is_default_role"`
	Permissions   []string `json:"permissions"`
}

type requestRoleUpdate struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	IsDefaultRole bool     `json:"is_default_role"`
	Permissions   []string `json:"permissions"`
}

type requestRoleScopeCreate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type responseRoleGet struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Role    RoleResource `json:"role"`
}

type responseRoleCreate struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Role    RoleResource `json:"role"`
}

type responseRoleUpdate struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Role    RoleResource `json:"role"`
}

type responseRoleDelete struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type responseRolePermissions struct {
	Permissions []string `json:"permissions"`
}

type PermissionResource struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Key         string `json:"key"`
}

type requestPermissionCreate struct {
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string `json:"description"`
}

type requestPermissionUpdate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type responsePermissionGet struct {
	Code       string             `json:"code"`
	Message    string             `json:"message"`
	Permission PermissionResource `json:"permission"`
}

type responsePermissionCreate struct {
	Code       string             `json:"code"`
	Message    string             `json:"message"`
	Permission PermissionResource `json:"permission"`
}

type responsePermissionUpdate struct {
	Code       string             `json:"code"`
	Message    string             `json:"message"`
	Permission PermissionResource `json:"permission"`
}

type responseRoleScopesGet struct {
	Code    string             `json:"code"`
	Message string             `json:"message"`
	Scopes  []ApiScopeResource `json:"scopes"`
}

type responseApiDelete struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type responseApiScopeDelete struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type responsePermissionDelete struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type responseRoleScopeDelete struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
