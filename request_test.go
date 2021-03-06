package fastly

import (
	"net/url"
	"strings"
	"testing"
)

func TestClient_RawRequest(t *testing.T) {
	validAPIHosts := []string{
		"https://api.fastly.com",
		"https://api.fastly.com/",
	}
	purgeAPIPaths := []string{
		"/service/myservice/purge/",
		"service/myservice/purge/",
	}
	cacheKeys := []string{
		"/",
		"text//text",
		"$-_.+!*'(),,;/?:@=&\"<>#%{}|\\^~[]`",
	}
	c := &Client{}
	for _, h := range validAPIHosts {
		c.url = &url.URL{Path: h}
		for _, p := range purgeAPIPaths {
			for _, k := range cacheKeys {
				r, err := c.RawRequest("GET", p+k, nil)
				// Cannot test results for success if we get an error
				if err != nil {
					t.Fatal("Could not make RawRequest for ", h, p, k)
				}
				t.Log("Path returned: ", r.URL.Path)
				pk := p + k
				if p[0] != '/' {
					pk = "/" + pk
				}
				t.Log("Path expected: ", pk)
				// Insure we don't get a path starting with an extra slash
				// e.g. //service/myservice/purge/mykey
				if r.URL.Path[1] == '/' {
					t.Fatal("Host and APIPath were joined incorrectly.")
				}
				// Insure the cache key isn't altered
				if strings.Index(r.URL.Path, p+k) == -1 {
					t.Fatal("RawRequest altered the cache key")
				}
			}
		}
	}
}
