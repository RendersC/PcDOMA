package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

const apiPrefix = "/api/v1"

func To(target string) gin.HandlerFunc {
	targetURL, err := url.Parse(target)
	if err != nil {
		panic("invalid proxy target: " + target)
	}

	rp := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			req.Host = targetURL.Host
			req.Header.Set("X-Forwarded-Host", req.Host)
			// Strip /api/v1 prefix: /api/v1/auth/register → /auth/register
			path := req.URL.Path
			if strings.HasPrefix(path, apiPrefix) {
				path = path[len(apiPrefix):]
				if path == "" {
					path = "/"
				}
			}
			req.URL.Path = path
			req.URL.RawPath = ""
		},
		ModifyResponse: func(resp *http.Response) error {
			resp.Header.Set("X-Served-By", "PCDoma Gateway")
			return nil
		},
	}

	return func(c *gin.Context) {
		rp.ServeHTTP(c.Writer, c.Request)
	}
}
