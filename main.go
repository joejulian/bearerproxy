package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

const (
	// Defaults
	defaultListenPort = "8080"
	defaultTargetURL  = "http://localhost:80/"

	// Environment strings
	envPort        = "PORT"
	envTargetURL   = "TARGET_URL"
	envBearerToken = "BEARER_TOKEN"
)

type proxy struct {
	listenPort  string
	targetURL   string
	bearerToken string
}

// Get env var or default
func getEnv(key string, defaultValue *string) *string {
	if value, ok := os.LookupEnv(key); ok {
		return &value
	}
	return defaultValue
}

// Get the port to listen on
func getListenPort() string {
	d := defaultListenPort
	port := getEnv(envPort, &d)
	return *port
}

// Get the url to proxy to
func getTargetURL() string {
	d := defaultTargetURL
	url := getEnv(envTargetURL, &d)
	return *url
}

// Get token to use for bearer authentication
func getBearerToken() *string {
	url := getEnv(envBearerToken, nil)
	return url
}

// log the setup
func (p *proxy) logSetup() {
	p.listenPort = getListenPort()
	p.targetURL = getTargetURL()
	token := getBearerToken()
	if token == nil {
		log.Fatal("BEARER_TOKEN environment variable must be set")
	}
	p.bearerToken = *token

	log.Printf("listening on port %s\n", p.listenPort)
	log.Printf("proxying for %s\n", p.targetURL)
}

func (p *proxy) serveReverseProxy(res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, err := url.Parse(p.targetURL)
	if err != nil {
		log.Fatal("error parsing url", err)
	}

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Header.Set("Authorization", "Bearer "+p.bearerToken)
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

func (p *proxy) handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	p.serveReverseProxy(res, req)
}

func main() {
	proxy := new(proxy)
	proxy.logSetup()

	// start server
	http.HandleFunc("/", proxy.handleRequestAndRedirect)
	if err := http.ListenAndServe(":"+getListenPort(), nil); err != nil {
		log.Fatal(err)
	}
}

