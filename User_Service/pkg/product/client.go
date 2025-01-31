package product

import (
	"net/http"
	"time"
)

type ClientProduct struct {
	Client  *http.Client
	BaseUrl string
	ApiKey  string
}

func NewClientProducts() *ClientProduct {

	return &ClientProduct{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		BaseUrl: "http://productservice:8083",
		ApiKey:  "lwehvowhvowvhwovwfwefwefwefw",
	}
}
