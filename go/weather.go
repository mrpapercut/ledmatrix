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

type Weather struct {
	ApiKey         string
	Location       int
	EndpointPrefix string
}

var weatherLock = &sync.Mutex{}
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

func (w *Weather) GetIcon(icon int) string {
	icons := map[int]string{
		1:  "sun",
		2:  "sun",
		3:  "sun",
		4:  "cloudy",
		5:  "cloudy",
		6:  "cloudy",
		7:  "clouds",
		8:  "clouds",
		9:  "",
		10: "",
		11: "clouds",
		12: "rain",
		13: "rain",
		14: "rain",
		15: "lightning",
		16: "lightning",
		17: "lightning",
		18: "rain",
		19: "snow",
		20: "snow",
		21: "snow",
		22: "snow",
		23: "snow",
		24: "snow",
		25: "snow",
		26: "rain",
		27: "",
		28: "",
		29: "snow",
		30: "thermometer_high",
		31: "thermometer_low",
		32: "windy",
		33: "moon",
		34: "moon",
		35: "moon",
		36: "clouds",
		37: "clouds",
		38: "clouds",
		39: "rain",
		40: "rain",
		41: "lightning",
		42: "lightning",
		43: "snow",
		44: "snow",
	}

	if iconName, ok := icons[icon]; ok {
		return iconName
	}

	return ""
}
