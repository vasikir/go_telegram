package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	TOKEN   = "**"
	WEATHER = "**"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Moscow"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("London"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Tula"),
	),
)

func main() {
	bot, err := tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		msg := tgbotapi.NewMessage(chatID, update.Message.Text)

		switch update.Message.Text {
		case "open":
			msg.ReplyMarkup = numericKeyboard
		case "Moscow", "London", "Tula":
			var w Weather
			err := getWeather(&w, update.Message.Text)
			if err != nil {
				panic(err)
			}
			msg_weather := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", w.Current.TempC))
			if _, err := bot.Send(msg_weather); err != nil {
				log.Panic(err)
			}
		case "close":
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

type Weather struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		Uv         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
		AirQuality struct {
			Co           float64 `json:"co"`
			No2          float64 `json:"no2"`
			O3           float64 `json:"o3"`
			So2          float64 `json:"so2"`
			Pm25         float64 `json:"pm2_5"`
			Pm10         float64 `json:"pm10"`
			UsEpaIndex   int     `json:"us-epa-index"`
			GbDefraIndex int     `json:"gb-defra-index"`
		} `json:"air_quality"`
	} `json:"current"`
}

func getWeather(w *Weather, city string) error {
	url := "http://api.weatherapi.com/v1/current.json?q=" + city + "&aqi=yes&key=" + WEATHER

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	jsonData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(jsonData), w)

	if err != nil {
		return err
	}

	return err
}
