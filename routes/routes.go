// Package routes. DO NOT EDIT

package routes

import (
        "github.com/gin-gonic/gin"
        "net/http"
        "github.com/gin-contrib/gzip"
        cors "github.com/rs/cors/wrapper/gin"
		"github.com/thanhhh/gin-gonic-realip"
    )

type Route struct {
    Path   string
    Route  func(*gin.RouterGroup)
}

var DefaultRoutes = []Route{
    {
        Path:  "/hello",
        Route: HelloRoute,
    },
}


func New() *gin.Engine {
	r := gin.New()
	initRoute(r)
	r.Use(realip.RealIP())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(cors.Default())


	api := r.Group("/api")
    for _, route := range DefaultRoutes {
        route.Route(api)
    }
	return r
}

func initRoute(r *gin.Engine) {
	_ = r.SetTrustedProxies(nil)
	r.RedirectTrailingSlash = false
	r.HandleMethodNotAllowed = true

	r.NoRoute(func(c *gin.Context) {
		SendErrorResponse(c, http.StatusNotFound, c.Request.RequestURI+" not found")
	})

	r.NoMethod(func(c *gin.Context) {
		SendErrorResponse(c, http.StatusMethodNotAllowed, c.Request.Method+" is not allowed here")
	})
}

func InitGin() {
	gin.DisableConsoleColor()
}


type Response struct {
	StatusCode int            `json:"-"`
	Success    bool           `json:"success"`
	Message    string         `json:"message,omitempty"`
	Data       map[string]any `json:"data,omitempty"`
}

func (response *Response) SendResponse(c *gin.Context) {
	c.AbortWithStatusJSON(response.StatusCode, response)
}

func SendErrorResponse(c *gin.Context, status int, message string) {
	response := &Response{
		StatusCode: status,
		Success:    false,
		Message:    message,
	}
	response.SendResponse(c)
}

func SendResponseData(c *gin.Context, data gin.H) {
	response := &Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       data,
	}
	response.SendResponse(c)
}