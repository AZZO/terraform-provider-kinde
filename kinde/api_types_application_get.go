package kinde

type ApiResponseApplicationGet struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Application struct {
		Id           string `json:"id"`
		Name         string `json:"name"`
		Type         string `json:"type"`
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}
}
