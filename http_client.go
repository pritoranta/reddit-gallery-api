package main

import "net/http"

const userAgent = "web:pritoranta-gallery-api:1.1.0 (by /u/AstronomerFit3983)"

type defaultHeaderRoundTripper struct {
	rt http.RoundTripper
}

func (dhrt *defaultHeaderRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", userAgent)
	return dhrt.rt.RoundTrip(req)
}

var HttpClient = &http.Client{
	Transport: &defaultHeaderRoundTripper{
		rt: http.DefaultTransport,
	}}
