package kinde

type ApiRequestApplicationCreate struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ApiResponseApplicationCreate struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Application struct {
		Id           string `json:"id"`
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client"`
	}
}
