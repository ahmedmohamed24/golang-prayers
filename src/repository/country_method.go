package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ahmedmohamed24/azan/config"
	"github.com/ahmedmohamed24/azan/models"
	"github.com/hablullah/go-prayer"
	"gorm.io/gorm"
)

func CountryDefaultMethod(country string, app *config.APP) prayer.CalculationMethod {
	cacheKey := fmt.Sprintf("country_method_%v", country)
	cachedValue, err := app.REDIS_DB.Get(cacheKey).Result()
	if err == nil {
		method, err := strconv.Atoi(cachedValue)
		if err != nil {
			fmt.Println("restoring from cached version")
			return prayer.CalculationMethod(method)
		}
	}
	// fetch the method from the db and cache it
	countryMehtod := models.CountryMethod{}
	result := app.DB.Where("country_name = ?", country).First(&countryMehtod)
	if result.Error != nil {
		fmt.Println(fmt.Errorf("%s", result.Error))
		return prayer.MWL
	}

	calculationMehtodsAdapter := map[string]prayer.CalculationMethod{
		"EGYPT":        prayer.Egypt,
		"RUSSIA":       prayer.MWL,
		"MOONSIGHTING": prayer.MWL,
		"SINGAPORE":    prayer.MUIS,
		"ISNA":         prayer.ISNA,
		"TURKEY":       prayer.MWL,
		"FRANCE":       prayer.France15,
		"QATAR":        prayer.Gulf,
		"KUWAIT":       prayer.Gulf,
		"MAKKAH":       prayer.UmmAlQura,
	}

	if method, ok := calculationMehtodsAdapter[strings.ToUpper(countryMehtod.CountryName)]; ok {
		status := app.REDIS_DB.Set(cacheKey, strconv.Itoa(int(method)), time.Hour*24)
		if status.Val() != "OK" {
			fmt.Println(fmt.Errorf("%s", result.Error))
		}
		return method
	}
	app.REDIS_DB.Set(cacheKey, prayer.MWL, time.Hour*24)
	return prayer.MWL

}

func CountryMethodSeeding(db *gorm.DB) {
	seederName := "country_methods_json_file"
	seederResult := db.Where("name = ?", seederName).First(&models.Seeding{})
	if seederResult.Error != gorm.ErrRecordNotFound {
		return
	}
	seedingResult := db.Create(&models.Seeding{
		Name: seederName,
	})
	if seedingResult.Error != nil {
		panic(seedingResult)
	}
	//checks if there are data, stop seeding
	existsResult := db.First(&models.CountryMethod{})
	if existsResult.RowsAffected > 0 {
		return
	}
	file, err := os.Open("./config/country-method.json")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	var countryMethodsJson string
	for scanner.Scan() {
		countryMethodsJson += scanner.Text()
	}
	countryMethods := map[string]string{}

	err = json.Unmarshal([]byte(countryMethodsJson), &countryMethods)
	if err != nil {
		panic(err)
	}
	var countryMethodsModels []*models.CountryMethod
	for countryName, countryMethod := range countryMethods {
		countryMethodsModels = append(countryMethodsModels, &models.CountryMethod{
			CountryName: countryName,
			MethodName:  countryMethod,
		})

	}
	result := db.Create(countryMethodsModels)
	if result.Error != nil {
		panic(result.Error)
	}

	fmt.Println(result.RowsAffected)

}
