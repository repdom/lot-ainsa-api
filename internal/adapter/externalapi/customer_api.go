package externalapi

import (
	"be-lotsanmateo-api/internal/adapter/externalapi/model"
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
	"time"
)

type CustomerAPI struct {
	client  *http.Client
	baseURL url.URL
}

func (api *CustomerAPI) ExistCustomer(jwt, user, lang, document, documentType string) (error, error) {
	log.Println("ExistCustomer")

	change := url.URL{
		Scheme: api.baseURL.Scheme,
		Host:   api.baseURL.Host,
		Path:   api.baseURL.Path,
	}

	change.Path = path.Join(change.Path, "exist")
	params := change.Query()
	params.Add("document", document)
	params.Add("documentType", documentType)
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

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error reading response: %w", err)

	}

	log.Println("responser exist customer ", resp.Status)

	switch strings.TrimSpace(resp.Status) {
	case strconv.Itoa(http.StatusBadRequest):
		return fmt.Errorf("el cliente con el documento (%s) ya esta creado; code: %d", document, resp.StatusCode), nil
	case strconv.Itoa(http.StatusInternalServerError):
		return nil, fmt.Errorf("internal server error: %d", resp.StatusCode)
	default:
		return nil, nil
	}

}

func (api *CustomerAPI) CreateCustomer(jwt, user, lang string, customer model.Customer) (error, error) {
	log.Println("CreateCustomer")

	change := url.URL{
		Scheme: api.baseURL.Scheme,
		Host:   api.baseURL.Host,
		Path:   api.baseURL.Path,
	}

	log.Println(change.String())
	jsonData, err := json.Marshal(customer)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error marshalling json: %w", err)
	}
	log.Println(string(jsonData))
	req, err := http.NewRequest("POST", change.String(), bytes.NewBuffer(jsonData))
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		log.Println("Error al leer la respuesta:", err)
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	detail, ok := result["detail"].(string)
	if !ok {
		detail = "Internal Server Error"
	}

	log.Println("responser create customer ", detail)

	switch strconv.Itoa(resp.StatusCode) {
	case strconv.Itoa(http.StatusBadRequest):
		return fmt.Errorf("el cliente no se pudo crear: %s, %d", detail, resp.StatusCode), nil
	case strconv.Itoa(http.StatusNotFound):
		return fmt.Errorf("el cliente no se pudo crear: %s, %d", detail, resp.StatusCode), nil
	case strconv.Itoa(http.StatusInternalServerError):
		return nil, fmt.Errorf("internal server error: %s %d", detail, resp.StatusCode)
	default:
		return nil, nil
	}

}

func NewCustomerClient(baseURL string) *CustomerAPI {
	service, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	invokeService := url.URL{
		Scheme: service.Scheme,
		Host:   service.Host,
		Path:   "/api/v1/customers",
	}

	return &CustomerAPI{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		baseURL: invokeService,
	}
}
