// Package Types. DO NOT TOUCH
package types

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	StatusCode int            `json:"-"`
	Success    bool           `json:"success"`
	Message    string         `json:"message,omitempty"`
	Data       map[string]any `json:"data,omitempty"`
}

func (response *Response) SendResponse(c *gin.Context) {
	c.AbortWithStatusJSON(response.StatusCode, response)
}

func SendResponseData(c *gin.Context, data gin.H) {
	response := &Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       data,
	}
	response.SendResponse(c)
}

func SendErrorResponse(c *gin.Context, status int, message string) {
	response := &Response{
		StatusCode: status,
		Success:    false,
		Message:    message,
	}
	response.SendResponse(c)
}

// Structs for API responses
type IPResponse struct {
    IP string `json:"ip"`
}

type LocationResponse struct {
    Status     string  `json:"status"`
    Country    string  `json:"country"`
    RegionName string  `json:"regionName"`
    City       string  `json:"city"`
    Lat        float64 `json:"lat"`
    Lon        float64 `json:"lon"`
}

type WeatherResponse struct {
    Weather []struct {
        Description string `json:"description"`
    } `json:"weather"`
    Main struct {
        Temp float64 `json:"temp"`
    } `json:"main"`
}