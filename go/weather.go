package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type UnitValue struct {
	Value float64 `json:"Value"`
	Unit  string  `json:"Unit"`
}

type Condition struct {
	Time        int    `json:"EpochTime"`
	WeatherText string `json:"WeatherText"`
	WeatherIcon int    `json:"WeatherIcon"`
	Temperature struct {
		Metric UnitValue `json:"Metric"`
	} `json:"Temperature"`
	PrecipitationSummary struct {
		Precipitation struct {
			Metric UnitValue `json:"Metric"`
		} `json:"Precipitation"`
	} `json:"PrecipitationSummary"`
}

type CurrentConditions []Condition

type DayPart struct {
	Icon            int       `json:"Icon"`
	Phrase          string    `json:"ShortPhrase"`
	RainProbability int       `json:"RainProbability"`
	Rain            UnitValue `json:"Rain"`
	Snow            UnitValue `json:"Snow"`
}

type DailyForecast struct {
	Date int `json:"EpochDate"`
	Sun  struct {
		Rise int `json:"EpochRise"`
		Set  int `json:"EpochSet"`
	} `json:"Sun"`
	Moon struct {
		Rise int `json:"EpochRise"`
		Set  int `json:"EpochSet"`
	} `json:"Moon"`
	Temperature struct {
		Minimum UnitValue `json:"Minimum"`
		Maximum UnitValue `json:"Maximum"`
	} `json:"Temperature"`
	Day   DayPart `json:"Day"`
	Night DayPart `json:"Night"`
}

type FiveDayForecast struct {
	DailyForecasts []DailyForecast `json:"DailyForecasts"`
}

var weatherLock = &sync.Mutex{}

type Weather struct {
	ApiKey         string
	Location       int
	EndpointPrefix string
}

var weatherInstance *Weather

func getWeatherInstance() *Weather {
	config := getConfig()

	if weatherInstance == nil {
		weatherLock.Lock()
		defer weatherLock.Unlock()

		if weatherInstance == nil {
			weatherInstance = &Weather{
				ApiKey:         config.Weather.ApiKey,
				Location:       config.Weather.Location,
				EndpointPrefix: "https://dataservice.accuweather.com",
			}
		}
	}

	return weatherInstance
}

func (w *Weather) GetCurrentConditions() {
	endpointCurrentConditions := fmt.Sprintf("%v/currentconditions/v1/%v?apikey=%v&details=true", w.EndpointPrefix, w.Location, w.ApiKey)

	response, err := http.Get(endpointCurrentConditions)
	if err != nil {
		fmt.Println("Error getting current conditions:", err)
		return
	}
	defer response.Body.Close()

	var currentConditions CurrentConditions

	err = json.NewDecoder(response.Body).Decode(&currentConditions)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	fmt.Println(currentConditions)
}

func (w *Weather) Get5DayForecast() {
	endpointFiveDayForecast := fmt.Sprintf("%v/forecasts/v1/daily/5day/%v?apikey=%v&details=true&metric=true", w.EndpointPrefix, w.Location, w.ApiKey)

	response, err := http.Get(endpointFiveDayForecast)
	if err != nil {
		fmt.Println("Error getting 5day forecast:", err)
		return
	}
	defer response.Body.Close()

	var fiveDayForecast FiveDayForecast

	err = json.NewDecoder(response.Body).Decode(&fiveDayForecast)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	fmt.Println(fiveDayForecast)
}
