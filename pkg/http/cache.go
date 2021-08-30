package http

import (
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/go-resty/resty/v2"
)

var cache *ristretto.Cache

func init() {
	c, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1024 * 10,
		MaxCost:     1024 * 1024 * 1024,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	cache = c
}

func cacheGet(url string) (Response, bool) {
	r, ok := cache.Get(url)
	if ok {
		v, ok := r.(Res)
		if ok {
			return v, true
		}

		cache.Del(url)
		cache.Wait()

		return nil, false
	}

	return nil, false
}

func cacheSet(url string, resp *resty.Response) Response {
	res := Res{
		url:  resp.Request.URL,
		code: resp.StatusCode(),
		body: resp.Body(),
	}
	if res.code < 300 {
		cache.SetWithTTL(url, res, int64(len(resp.Body())), time.Hour)
		cache.Wait()
	}
	return res
}

type Response interface {
	Body() []byte
	String() string
	StatusCode() int
	URL() string
}

type Res struct {
	code int
	url  string
	body []byte
}

func (r Res) StatusCode() int {
	return r.code
}

func (r Res) Body() []byte {
	return r.body
}

func (r Res) String() string {
	return string(r.body)
}

func (r Res) URL() string {
	return r.url
}
