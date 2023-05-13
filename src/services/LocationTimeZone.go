package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

const (
	GOOGLE_TIMEZONE = "https://maps.googleapis.com/maps/api/timezone/json"
)

type Timezone struct {
	Status     string `json:"status"`
	TimeZoneID string `json:"TimeZoneId"`
}

func LocationTimeZone(lat, lng float64) (*Timezone, error) {

	resp, err := http.Get(fmt.Sprintf("%s?key=%s&location=%v,%v&timestamp=1331161200", GOOGLE_TIMEZONE, os.Getenv("GOOGLE_MAPS_KEY"), lat, lng))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var timezone Timezone
	if err := json.NewDecoder(resp.Body).Decode(&timezone); err != nil {
		return nil, err
	}
	if timezone.Status != "OK" {
		return nil, errors.New("something went wrong during fetching timezone from our thirdparty")
	}
	return &timezone, nil

}
