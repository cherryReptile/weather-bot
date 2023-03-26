package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type WeatherResponse struct {
	Coord      Coord     `json:"coord"`
	Weather    []Weather `json:"weather"`
	Base       string    `json:"base"`
	Main       Main      `json:"main"`
	Visibility int       `json:"visibility"`
	Wind       Wind      `json:"wind"`
	Rain       Rain      `json:"rain"`
	Snow       Snow      `json:"snow"`
	Clouds     Clouds    `json:"clouds"`
	Dt         int       `json:"dt"`
	Sys        Sys       `json:"sys"`
	Timezone   int       `json:"timezone"`
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Cod        int       `json:"cod"`
}

type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Weather struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}

type Rain struct {
	H  float64 `json:"1h"`
	H3 float64 `json:"3h"`
}

type Snow struct {
	H  float64 `json:"1h"`
	H3 float64 `json:"3h"`
}

type Clouds struct {
	All int `json:"all"`
}

type Sys struct {
	Type    int    `json:"type"`
	Id      int    `json:"id"`
	Country string `json:"country"`
	Sunrise int    `json:"sunrise"`
	Sunset  int    `json:"sunset"`
}

func (w *WeatherResponse) GetWeather(lat float64, lng float64) error {
	// Create a new HTTP client
	client := http.Client{}

	// Build the request URL
	u := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%v&lon=%v&appid=%s", lat, lng, os.Getenv("OPENWEATHERMAP_API_KEY"))

	// Send the request
	resp, err := client.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(w)
	if err != nil {
		return err
	}

	return nil
}

func (w *WeatherResponse) GetTemperature() string {
	return fmt.Sprintf("Temperature 🌡:\n\t\tIn Celcius: %.2f°C\n\t\tIn Fahrenheit: %.2f°F\n\t\tIn Kelvin: %.2fK", w.Main.Temp-273.15, (w.Main.Temp-273.15)*9/5+32, w.Main.Temp)
}

func (w *WeatherResponse) GetWeatherStat() string {
	weather := "Weather 🌤:"
	s := "\t\tState: "
	d := "\t\tDescription: "
	for _, v := range w.Weather {
		s += v.Main
		d += v.Description
	}
	return fmt.Sprintf("%s\n%s\n%s", weather, s, d)
}

func (w *WeatherResponse) GetRainStat() string {
	var l1 string
	var l3 string
	res := "Rain 🌧:"
	if w.Rain.H != 0 {
		l1 = fmt.Sprintf("last 1h: %v mm", w.Rain.H)
	}
	if w.Rain.H != 0 {
		l3 = fmt.Sprintf("last 3h: %v mm", w.Rain.H3)
	}
	if l1 != "" {
		res += fmt.Sprintf("\n\t\t%s", l1)
	}
	if l3 != "" {
		res += fmt.Sprintf("\n\t\t%s", l3)
	}
	if l1 == "" && l3 == "" {
		res = ""
	}

	return res
}

func (w *WeatherResponse) GetSnowStat() string {
	var l1 string
	var l3 string
	res := "Snow ❄️: "
	if w.Rain.H != 0 {
		l1 = fmt.Sprintf("last 1h: %v mm", w.Rain.H)
	}
	if w.Rain.H != 0 {
		l3 = fmt.Sprintf("last 3h: %v mm", w.Rain.H3)
	}
	if l1 != "" {
		res += fmt.Sprintf("\n\t\t%s", l1)
	}
	if l3 != "" {
		res += fmt.Sprintf("\n\t\t%s", l3)
	}
	if l1 == "" && l3 == "" {
		res = ""
	}

	return res
}

func (w *WeatherResponse) GetWindStat() string {
	return fmt.Sprintf("Wind 🌬:\n\t\tSpeed: %.2f\n\t\tDirection,Degrees: %v", w.Wind.Speed, w.Wind.Deg)
}

func (w *WeatherResponse) GetCloudsStat() string {
	return fmt.Sprintf("Clouds ☁️: %v%s", w.Clouds.All, "%")
}

func (w *WeatherResponse) GetHumidity() string {
	return fmt.Sprintf("Humidity 💧: %v%s", w.Main.Humidity, "%")
}

func (w *WeatherResponse) GetVisibility() string {
	return fmt.Sprintf("Visibility 👀: %vm", w.Visibility)
}

func (w *WeatherResponse) GetPressureStat() string {
	return fmt.Sprintf("Pressure 🗿: %v hPa", w.Main.Pressure)
}
func (w *WeatherResponse) GetCountryAndCity() string {
	return fmt.Sprintf("Country, city: %s, %s", w.Sys.Country, w.Name)
}

func (w *WeatherResponse) GetInfo() string {
	var rain string
	var snow string
	t := w.GetTemperature() + "\n\n"
	wstat := w.GetWeatherStat() + "\n\n"
	wind := w.GetWindStat() + "\n\n"
	cloud := w.GetCloudsStat() + "\n\n"
	hum := w.GetHumidity() + "\n\n"
	vis := w.GetVisibility() + "\n\n"
	pres := w.GetPressureStat() + "\n\n"

	if w.GetRainStat() != "" {
		rain = w.GetRainStat() + "\n\n"
	}
	if w.GetSnowStat() != "" {
		snow = w.GetSnowStat() + "\n\n"
	}

	cc := w.GetCountryAndCity()
	res := fmt.Sprint(t, wstat, wind, cloud, hum, vis, pres, rain, snow, cc)
	return res
}
