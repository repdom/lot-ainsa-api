package reservation

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

func NewReservationAPI(baseURL string) *API {
	service, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	invokeService := url.URL{
		Scheme: service.Scheme,
		Host:   service.Host,
		Path:   "/api/v1/reservations",
	}

	return &API{
		utility: externalapi.NewUtilityAPI(&http.Client{
			Timeout: 15 * time.Second,
		}, invokeService),
	}
}

func (api *API) GetReservation(jwt, user, lang string, id int) (*model.ReservationDomain, error) {
	log.Println("Loading Reservation")

	urlStr := api.utility.BuildURLWithID("id", id)

	req, err := api.utility.BuildRequestGet(urlStr, jwt, user, lang)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var result model.ReservationDomain
	if err := api.utility.DoRequest(req, &result); err != nil {
		log.Println(err)
		return nil, err
	}

	return &result, nil
}
