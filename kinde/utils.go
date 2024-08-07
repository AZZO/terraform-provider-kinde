package kinde

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

func VerifyToken(m2mToken *oauth2.Token, tokenIssuer string, audience string) error {
	jwksURL := fmt.Sprintf("%v/.well-known/jwks", tokenIssuer)

	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})

	if err != nil {
		return err
	}

	parsedToken, err := jwt.Parse(m2mToken.AccessToken, jwks.Keyfunc,
		jwt.WithValidMethods([]string{"RS256"}), //verifying the signing algorythm
		jwt.WithIssuer(tokenIssuer),             //verifying the token issuer
		jwt.WithAudience(audience))              //verifying that the token is for correct audience

	if err != nil {
		return fmt.Errorf("error verifying token %v", err)
	}

	fmt.Printf("m2m token %#v", parsedToken.Claims)

	return nil
}

func GetJson(client *http.Client, url string, result interface{}) error {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	response, err := client.Do(req)

	if err != nil {
		return err
	} else if !HttpStatusCodeIsSuccess(response.StatusCode) {
		return fmt.Errorf("received status code %v for url %v", response.StatusCode, url)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}

	return nil
}

func PostJson(client *http.Client, url string, request interface{}, result interface{}) error {

	reqJson, err := json.Marshal(request)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	response, err := client.Do(req)

	if err != nil {
		return err
	} else if !HttpStatusCodeIsSuccess(response.StatusCode) {
		return fmt.Errorf("received status code %v for url %v", response.StatusCode, url)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}

	return nil
}

func PatchJson(client *http.Client, url string, request interface{}) error {

	reqJson, err := json.Marshal(request)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(reqJson))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	response, err := client.Do(req)

	if err != nil {
		return err
	} else if !HttpStatusCodeIsSuccess(response.StatusCode) {
		return fmt.Errorf("received status code %v for url %v", response.StatusCode, url)
	}

	defer response.Body.Close()

	return nil
}

func Delete(client *http.Client, url string) error {

	req, err := http.NewRequest("DELETE", url, nil)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	response, err := client.Do(req)

	if err != nil {
		return err
	} else if !HttpStatusCodeIsSuccess(response.StatusCode) {
		return fmt.Errorf("received status code %v for url %v", response.StatusCode, url)
	}

	defer response.Body.Close()

	return nil
}

func HttpStatusCodeIsSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
