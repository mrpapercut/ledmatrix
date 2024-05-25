package types

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
