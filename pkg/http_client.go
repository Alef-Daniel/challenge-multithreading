package pkg

import (
	"net/http"
	"time"
)

type Client struct {
	Http *http.Client
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		Http: &http.Client{
			Timeout: timeout,
		},
	}
}
