package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

//Структура ответа от Api яндекса погоды
type Weather struct {
	//Время сервера в формате Unixtime
	Now      int64  `json:"now"`
	//Объект информации о населенном пункте
	Info Info `json:"info"`
	//Объект фактической информации о погоде
	Fact Fact  `json:"fact"`
	//Объект прогнозной информации о погоде
	Forecasts Forecasts `json:"forecasts"`
}

type Info struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
	Offset int64 `json:"offset"`
	Name string `json:"name"`
	Url string `json:"url"`
}

type Fact struct {
	Temp       float32 `json:"temp"`
	FeelsLike  float32 `json:"feels_like"`
	TempWater float32 `json:"temp_water"`
	WindSpeed   float32 `json:"wind_speed"`
	WindGust   float32 `json:"wind_gust"`
	PressureMm float32 `json:"pressure_mm"`
}

type Forecasts struct {
	Sunrise string `json:"sunrise"`
	Sunset string `json:"sunset"`
	Parts Parts  `json:"parts"`
}

type Parts struct {
	Night Part `json:"night"`
	Morning Part `json:"morning"`
	Day Part `json:"day"`
	Evening Part `json:"evening"`
}

type Part struct {
	TempMin    float32  `json:"temp_min"`
	TempMax    float32  `json:"temp_max"`
	FeelsLike float32  `json:"feels_like"`
	WindSpeed float32 `json:"wind_speed"`
}

type Error struct {
	Error string
}

type InfoResponse struct {
	Text string
}
func main() {
	e := echo.New()
	yandexApiKey := "6dfedf6f-1d92-4092-b8b7-46865ff715be"

	e.GET("/getWeather", func(c echo.Context) error {
		lat, errLat := strconv.ParseFloat(c.QueryParam("lat"), 64)
		lon, errLon := strconv.ParseFloat(c.QueryParam("lon"), 64)

		if errLat == nil && errLon ==nil {
			weather,_ := getWeather(yandexApiKey, "ru_RU", lat, lon)
				return c.JSON(http.StatusOK, InfoResponse{Text: getUserFriendlyTempInfo(weather.Fact.Temp)})
		} else {
			err := Error{"Ошибка!"}
			return c.JSON(http.StatusOK, err)
		}
	})
	e.Logger.Fatal(e.Start(":1323"))
}


func getWeather(apiKey string, lang string, lat , lon float64) (Weather, error) {
	url := fmt.Sprintf("https://api.weather.yandex.ru/v2/forecast?lat=%f&lon=%f&lang=%v", lat, lon, lang)

	client := http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Weather{}, err
	}
	request.Header.Add("X-Yandex-API-Key", apiKey)
	resp, err := client.Do(request)
	if err != nil {
		return Weather{}, err
	}
	defer resp.Body.Close()

	var weather Weather

	err = json.NewDecoder(resp.Body).Decode(&weather)
	return weather, nil
}

func getUserFriendlyTempInfo( temp float32) string {
	return fmt.Sprintf("Фактическая температура %.1f по цельсию", temp)
}