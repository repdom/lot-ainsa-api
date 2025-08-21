package financing

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

func (api *API) LoadFinancing(jwt, user, lang string, financingId int) (*model.FinancingDomain, error) {
	log.Println("Loading down payment")

	urlStr := api.utility.BuildURLWithID("id", financingId)

	req, err := api.utility.BuildRequestGet(urlStr, jwt, user, lang)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var result model.FinancingDomain
	if err := api.utility.DoRequest(req, &result); err != nil {
		log.Println(err)
		return nil, err
	}

	return &result, nil

}

func (api *API) PatchFinancing(jwt, user, lang string, financingId int, financingPath model.FinancingDomain) (*model.FinancingDomain, error) {
	log.Println("Loading financing")

	urlStr := api.utility.BuildURLParameter("id", financingId)

	req, err := api.utility.BuildRequestPatch(urlStr, jwt, user, lang, financingPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var result model.FinancingDomain
	if err := api.utility.DoRequest(req, &result); err != nil {
		log.Println(err)
		return nil, err
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
		Path:   "/api/v1/financings",
	}

	return &API{
		utility: externalapi.NewUtilityAPI(&http.Client{
			Timeout: 15 * time.Second,
		}, invokeService),
	}
}
