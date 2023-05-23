package opa2

import (
	"context"
	"net/http"
	"strings"
)

// config holds config to pass the plug-in
type Config struct{}

func CreateConfig() *Config {
	return &Config{}
}

type Opa2 struct {
	next http.Handler
	name string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &Opa2{
		next: next,
		name: name,
	}, nil

}

func (u *Opa2) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	var requestAgents []string = req.Header["User-Agent"]
	var agentListString string = strings.Join(requestAgents[:], ",")

	var allow bool = false // DEFAULT to NOT allow

	blockedAgents := [4]string{"yandex", "petal", "bot", "postman"}
	for _, agent := range blockedAgents {
		if strings.Contains(strings.ToLower(agentListString), agent) {
			rw.WriteHeader(403) // FORBIDDEN
		}
	}

	knownAgents := [8]string{"chrome", "mozilla", "edg", "applewebkit", "safari", "appspider", "bingbot", "googlebot"}
	for _, agent := range knownAgents {
		if strings.Contains(strings.ToLower(agentListString), agent) {
			allow = true
			break // We found < 1 leave loop
		}
	}

	if allow {
		u.next.ServeHTTP(rw, req) // continue with next middleware
	} else {
		rw.WriteHeader(403) // FORBIDDEN
	}
}
