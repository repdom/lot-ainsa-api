package company

import (
	"be-lotsanmateo-api/internal/adapter/externalapi"
	"be-lotsanmateo-api/internal/adapter/externalapi/model"
	"log"
	"net/http"
	"net/url"
	"time"
)

type API struct {
	utility *externalapi.UtilityAPI
}

func NewCompanyAPI(baseURL string) *API {
	service, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	invokeService := url.URL{
		Scheme: service.Scheme,
		Host:   service.Host,
		Path:   "/api/v1/company-configurations",
	}

	return &API{
		utility: externalapi.NewUtilityAPI(&http.Client{
			Timeout: 15 * time.Second,
		}, invokeService),
	}
}

func (api *API) GetActiveCompanyConfiguration(jwt, user, lang string) (*model.CompanyConfigurationDomain, error) {
	log.Println("Loading active company configuration")

	req, err := api.utility.BuildRequestGet(api.utility.BuildURLDefault(), jwt, user, lang)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var result model.CompanyConfigurationDomain
	if err := api.utility.DoRequest(req, &result); err != nil {
		log.Println(err)
		return nil, err
	}

	return &result, nil
}
