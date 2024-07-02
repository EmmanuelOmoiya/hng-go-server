package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HelloRoute(router *gin.RouterGroup) {
	hello := router.Group("/hello")
	{
		hello.GET(
			"",
			GetDetails,
		)
	}
}

type IPResponse struct {
	IP string `json:"ip"`
}

type LocationResponse struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	AS          string  `json:"as"`
	Query       string  `json:"query"`
}

type WeatherResponse struct {
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func GetDetails(c *gin.Context) {
	visitorName := c.Query("visitor_name")
	if visitorName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "visitor_name query is required",
		})
		return
	}

	ipResp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting IP address"})
		return
	}
	defer ipResp.Body.Close()
	ipBody, _ := ioutil.ReadAll(ipResp.Body)

	var ip IPResponse
	if err := json.Unmarshal(ipBody, &ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing IP response"})
		return
	}

	locationURL := fmt.Sprintf("http://ip-api.com/json/%s", ip.IP)
	locationResp, err := http.Get(locationURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting location"})
		return
	}
	defer locationResp.Body.Close()
	locationBody, _ := ioutil.ReadAll(locationResp.Body)

	var location LocationResponse
	if err := json.Unmarshal(locationBody, &location); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing location response"})
		return
	}

	apiKey := "efe8e4cc101b121d2a17433bb54e0e42"
	weatherURL := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&units=metric&appid=%s", location.Lat, location.Lon, apiKey)
	weatherResp, err := http.Get(weatherURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting weather"})
		return
	}
	defer weatherResp.Body.Close()
	weatherBody, _ := ioutil.ReadAll(weatherResp.Body)

	var weather WeatherResponse
	if err := json.Unmarshal(weatherBody, &weather); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing weather response"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"client_ip": ip.IP,
		"location":  location.City,
		"greeting":  fmt.Sprintf("Hello, %s! The temperature is %.2f degrees Celsius in %s", visitorName, weather.Main.Temp, location.City),
	})
}
