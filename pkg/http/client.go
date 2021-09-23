package http

import (
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

var client = resty.New()

func Get(u *url.URL) (Response, error) {
	if forbiddenURL(u) {
		return nil, errors.New("禁止发送请求到域名 " + u.Host)
	}

	resp, ok := cacheGet(u.String())
	if !ok || resp.StatusCode() >= 300 {
		res, err := client.R().Get(u.String())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get resource %s", u)
		}
		resp = cacheSet(u.String(), res)
	}

	return resp, nil
}

func forbiddenURL(u *url.URL) bool {
	switch u.Host {
	case "cdn.jsdelivr.net":
		return false
	}
	//return strings.HasSuffix(u.Host, ".github.io")
	return true
}
