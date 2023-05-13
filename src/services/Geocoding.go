package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

const (
	GOOGLE_MAPS = "https://maps.googleapis.com/maps/api/geocode/json"
)

type AddressResult struct {
	Country     string  `json:"country"`
	State       string  `json:"state"`
	City        string  `json:"city"`
	Street      string  `json:"street"`
	Postal_code string  `json:"postal_code"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	PlaceId     string  `json:"place_id"`
	CountryCode string  `json:"country_code"`
}
type Results struct {
	Results []Result `json:"results"`
	Status  string   `json:"status"`
}

type Result struct {
	AddressComponents []Address `json:"address_components"`
	FormattedAddress  string    `json:"formatted_address"`
	Geometry          Geometry  `json:"geometry"`
	PlaceId           string    `json:"place_id"`
}
type Geometry struct {
	Location Location `json:"location"`
}

type Address struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

func Geocoding(location Location) (AddressResult, error) {
	results, err := geocodeByGoogle(location)
	var addressResult AddressResult
	if err != nil {
		return addressResult, err
	}
	result := results.Results[0]
	addressResult.PlaceId = result.PlaceId
	addressResult.Lat = result.Geometry.Location.Lat
	addressResult.Lng = result.Geometry.Location.Lng
	for _, addressComponent := range result.AddressComponents {
		switch addressComponent.Types[0] {
		case "route":
			addressResult.Street = strings.Trim(addressResult.Street+" "+addressComponent.LongName, " ")
		case "street_number":
			addressResult.Street = strings.Trim(addressResult.Street+" "+addressComponent.LongName, " ")
		case "country":
			addressResult.Country = addressComponent.LongName
			addressResult.CountryCode = addressComponent.ShortName
		case "administrative_area_level_1":
			addressResult.State = addressComponent.LongName
		case "postal_code":
			addressResult.Postal_code = addressComponent.LongName
		case "administrative_area_level_2":
			addressResult.City = addressComponent.LongName
		}
	}
	return addressResult, nil

}
func geocodeByGoogle(location Location) (Results, error) {
	var results Results
	client := &http.Client{}
	key := os.Getenv("GOOGLE_MAPS_KEY")
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?latlng=%v,%v&sensor=false&key=%s", GOOGLE_MAPS, location.Lat, location.Lng, key), nil)
	if err != nil {
		return results, err
	}
	req.Header.Set("Accept-Language", "en-US")
	resp, err := client.Do(req)
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return results, err
	}
	if results.Status != "OK" {
		return results, errors.New("couldn't found the address on google locations")
	}
	return results, nil
}
