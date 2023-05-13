package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ahmedmohamed24/azan/config"
	"github.com/ahmedmohamed24/azan/repository"
	"github.com/ahmedmohamed24/azan/requests"
	"github.com/ahmedmohamed24/azan/resources"
	"github.com/ahmedmohamed24/azan/services"
	"github.com/gin-gonic/gin"
	"github.com/hablullah/go-prayer"
)

func FetchAzan(ctx *gin.Context, app *config.APP) {
	var azanRequest requests.AzanRequest
	if valid := azanRequest.Validate(ctx); !valid {
		return
	}
	//check if there is a cached key stored in the db, returns it
	cacheKey := fmt.Sprintf("%v-%v-%v", time.Now().Format("2006-01-02"), azanRequest.Lat, azanRequest.Lng)
	cachedValue, err := app.REDIS_DB.Get(cacheKey).Result()
	if err == nil {
		getCachedResponse(ctx, cachedValue)
		return
	}
	location, err := services.Geocoding(services.Location{Lat: azanRequest.Lat, Lng: azanRequest.Lng})
	if err != nil {
		resources.RespondWithError(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	timezone, err := services.LocationTimeZone(location.Lat, location.Lng)
	if err != nil {
		resources.RespondWithError(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	timezoneLocation, err := time.LoadLocation(timezone.TimeZoneID)
	if err != nil {
		resources.RespondWithError(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	date := time.Now().In(timezoneLocation)

	prayerTimes, err := prayer.Calculate(prayer.Config{
		Latitude:          location.Lat,
		Longitude:         location.Lng,
		CalculationMethod: repository.CountryDefaultMethod(location.Country, app),
	}, date)
	if err != nil {
		resources.RespondWithError(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	response := resources.AzanResource{
		Message:  "success",
		City:     location.City,
		TimeZone: timezone.TimeZoneID,
		Date:     date.Format("2006-01-02 15:04:05"),
		Fajr:     prayerTimes.Fajr.Format("15:04:05"),
		Sunrise:  prayerTimes.Sunrise.Format("15:04:05"),
		Zuhr:     prayerTimes.Zuhr.Format("15:04:05"),
		Asr:      prayerTimes.Asr.Format("15:04:05"),
		Maghrib:  prayerTimes.Maghrib.Format("15:04:05"),
		Isha:     prayerTimes.Isha.Format("15:04:05"),
	}
	cachedResult, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	app.REDIS_DB.Set(cacheKey, cachedResult, time.Hour*24)
	ctx.JSON(http.StatusOK, response)

}

func getCachedResponse(ctx *gin.Context, cachedValue string) {
	azanResource := resources.AzanResource{}

	err := json.Unmarshal([]byte(cachedValue), &azanResource)
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, azanResource)
}
