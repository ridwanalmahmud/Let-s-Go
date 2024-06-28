package main

import (
	"fmt"
	"io"
	"time"
	"encoding/json"
	"net/http"
	"os"
	"github.com/fatih/color"
)

type Weather struct{
	Location struct{
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
  Current struct{
		TempC float64 `json:"temp_c"`
		Condition struct{
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct{
		ForecastDay []struct{
			Hour []struct{
				TimeEpoch int64 `json:"time_epoch"`
				TempC float64 `json:"temp_c"`
				Condition struct{
			    Text string `json:"text"`
		    } `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	q := "London"
	if len(os.Args) >= 2{
		q = os.Args[1]
	}
	res, err := http.Get("https://api.weatherapi.com/v1/forecast.json?key=b87d62014b8f4f87ada55732241806&q=" + q + "&days=3&aqi=no&alerts=no")
	if err != nil{
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200{
		panic("Weather api not found!!")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil{
		panic(err)
	}
	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil{
		panic(err)
	}
	location, current, hours := weather.Location, weather.Current, weather.Forecast.ForecastDay[0].Hour
	fmt.Printf(
		"%s, %s, %.0f°C, %s\n",
		location.Name,
		location.Country,
		current.TempC,
		current.Condition.Text,
	)
	for _, hour := range hours{
		date := time.Unix(hour.TimeEpoch, 0)
		if date.Before(time.Now()){
			continue
		}
		message:= fmt.Sprintf(
			"%s, %.0f°C, %.0f%%, %s\n",
			date.Format("15:04"), 
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)
		if hour.ChanceOfRain < 40{
			fmt.Printf(message)
		}else{
			color.Red(message)
		}
	}
}
