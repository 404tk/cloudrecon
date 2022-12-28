package asset_check

import (
	"crypto/tls"
	"net/http"
	"time"
)

var Timeout = 15 * time.Second

func newClient() *http.Client {
	var redirectFunc = func(_ *http.Request, _ []*http.Request) error {
		// Tell the http client to not follow redirect
		return http.ErrUseLastResponse
	}

	transport := &http.Transport{
		MaxIdleConnsPerHost: -1,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS10,
		},
		DisableKeepAlives: true,
	}

	return &http.Client{
		Transport:     locFix{transport},
		Timeout:       Timeout,
		CheckRedirect: redirectFunc,
	}
}

type locFix struct {
	http.RoundTripper
}

func (lf locFix) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := lf.RoundTripper.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode == 301 || resp.StatusCode == 302 {
		if resp.Header.Get("Location") == "" {
			resp.Header.Set("Location", "https://404.tips")
		}
	}
	return resp, err
}

func requestBuilder(url, method string) *http.Request {
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:47.0) Gecko/20100101 Firefox/47.0")
	return req
}
