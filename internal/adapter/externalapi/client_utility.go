package externalapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

const (
	GET   = "GET"
	PATCH = "PATCH"
)

type UtilityAPI struct {
	client  *http.Client
	baseURL url.URL
}

func NewUtilityAPI(client *http.Client, baseURL url.URL) *UtilityAPI {
	return &UtilityAPI{
		client:  client,
		baseURL: baseURL,
	}
}

func (api *UtilityAPI) BuildURLDefault() string {
	return api.baseURL.String()
}

func (api *UtilityAPI) BuildURLWithID(endpoint string, id int) string {
	return api.BuildURL(endpoint, map[string]string{
		"id": strconv.Itoa(id),
	})
}

func (api *UtilityAPI) BuildURLParameter(name string, value int) string {
	u := url.URL{
		Scheme: api.baseURL.Scheme,
		Host:   api.baseURL.Host,
		Path:   api.baseURL.Path,
	}

	params := u.Query()
	params.Add(name, strconv.Itoa(value))
	u.RawQuery = params.Encode()
	return u.String()
}

func (api *UtilityAPI) BuildURL(endpoint string, queryParams map[string]string) string {
	u := url.URL{
		Scheme: api.baseURL.Scheme,
		Host:   api.baseURL.Host,
		Path:   path.Join(api.baseURL.Path, endpoint),
	}

	params := u.Query()
	for k, v := range queryParams {
		params.Add(k, v)
	}
	u.RawQuery = params.Encode()
	return u.String()
}

func (api *UtilityAPI) BuildRequestGet(urlStr, jwt, user, lang string) (*http.Request, error) {
	return api.BuildRequest(GET, urlStr, jwt, user, lang, nil)
}

func (api *UtilityAPI) BuildRequestPatch(urlStr, jwt, user, lang string, target interface{}) (*http.Request, error) {
	jsonData, err := json.Marshal(target)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error marshalling request json: %w", err)
	}
	log.Println(string(jsonData))

	return api.BuildRequest(PATCH, urlStr, jwt, user, lang, bytes.NewBuffer(jsonData))
}

func (api *UtilityAPI) BuildRequest(method, urlStr, jwt, user, lang string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-language", lang)
	req.Header.Set("x-user", user)

	if jwt != "" && jwt != "null" && jwt != "undefined" {
		if !strings.HasPrefix(jwt, "Bearer ") {
			req.Header.Set("Authorization", "Bearer "+jwt)
		} else {
			req.Header.Set("Authorization", jwt)
		}
	}

	return req, nil
}

func (api *UtilityAPI) DoRequest(req *http.Request, target interface{}) error {
	resp, err := api.client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading body response: %w", err)
	}

	log.Printf("Status: %s\nBody: %s", resp.Status, string(body))

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("error decoding json: %w", err)
	}

	return nil
}
