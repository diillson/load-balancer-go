package loadbalancer

import (
	"net/http"
	"net/url"
)

type Backend struct {
	URL         *url.URL
	ActiveConns int32
	Healthy     bool
}

func (b *Backend) CheckHealth() {
	resp, err := http.Get(b.URL.String() + "/health")
	if err != nil || resp.StatusCode != http.StatusOK {
		b.Healthy = false
		return
	}
	b.Healthy = true
}
