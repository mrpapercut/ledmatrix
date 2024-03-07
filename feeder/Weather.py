import requests

class Weather:
    def __init__(self, config, db) -> None:
        self.db = db

        self.accuweather_api_key = config.get("weather").get("apikey")
        self.location_key = config.get("weather").get("location")

        self.endpoint_prefix = "https://dataservice.accuweather.com/"
        self.endpoint_current_conditions = f"{self.endpoint_prefix}currentconditions/v1/{self.location_key}?apikey={self.accuweather_api_key}&details=true"
        self.endpoint_5day_forecast = f"{self.endpoint_prefix}forecasts/v1/daily/5day/{self.location_key}?apikey={self.accuweather_api_key}&details=true&metric=true"

    def get_current_conditions(self):
        r = requests.get(self.endpoint_current_conditions)
        json_response = r.json()[0]

        temperature = json_response['Temperature']['Metric']
        realFeelTemperature = json_response['RealFeelTemperature']['Metric']
        precipitation = json_response['PrecipitationSummary']['Precipitation']['Metric']

        weather = {
            "timestamp": json_response['EpochTime'],
            "temperature": f"{temperature['Value']}{temperature['Unit']}",
            "real_feel_temperature": f"{realFeelTemperature['Value']}{realFeelTemperature['Unit']} ({realFeelTemperature['Phrase']})",
            "precipitation": f"{precipitation['Value']}{precipitation['Unit']}"
        }

        self.db.insert_current_weather(weather)

    def get_5day_forecast(self):
        r = requests.get(self.endpoint_5day_forecast)
        json_response = r.json()

        print(json_response)
