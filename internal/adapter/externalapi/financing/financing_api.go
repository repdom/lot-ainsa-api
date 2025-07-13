package financing

import (
	"be-lotsanmateo-api/internal/adapter/externalapi/model/financing"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type API struct {
	client  *http.Client
	baseURL url.URL
}

func (api *API) LoadFinancing(jwt, user, lang string, financingId int) (*financing.Financings, error) {
	log.Println("Load financing")
	change := url.URL{
		Scheme: api.baseURL.Scheme,
		Host:   api.baseURL.Host,
		Path:   api.baseURL.Path,
	}
	params := change.Query()
	params.Add("id", strconv.Itoa(financingId))
	change.RawQuery = params.Encode()
	log.Println(change.String())

	req, err := http.NewRequest("GET", change.String(), nil)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("x-language", lang)
	req.Header.Set("x-user", user)

	if jwt != "" && jwt != "null" && jwt != "undefined" {
		if !strings.Contains(jwt, "Bearer") {
			req.Header.Set("Authorization", "Bearer "+jwt)
		} else {
			req.Header.Set("Authorization", jwt)
		}
	}

	resp, err := api.client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error making request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error reading response: %w", err)

	}

	var result financing.Financings
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}
	return &result, nil

}

func NewFinancingAPI(baseURL string) *API {
	service, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	invokeService := url.URL{
		Scheme: service.Scheme,
		Host:   service.Host,
		Path:   "/api/v1/financings/id",
	}

	return &API{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		baseURL: invokeService,
	}
}
