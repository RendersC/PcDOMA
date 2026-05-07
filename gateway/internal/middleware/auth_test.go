package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// newRouter wires a single protected route guarded by Required(). The final
// handler echoes the identity headers the middleware injects, so tests can
// assert both the status code and that the user context was propagated.
func newRouter(authURL string) *gin.Engine {
	m := NewAuthMiddleware(authURL, nil) // nil redis → always falls through to HTTP validate
	r := gin.New()
	r.GET("/protected", m.Required(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user_id": c.Request.Header.Get("X-User-ID"),
			"role":    c.Request.Header.Get("X-User-Role"),
		})
	})
	return r
}

// stubAuth returns an httptest server standing in for auth-service's
// /internal/validate-token endpoint, replying with the given JSON body.
func stubAuth(body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
}

func TestRequired_MissingHeader(t *testing.T) {
	r := newRouter("http://unused")
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "authorization header required")
}

func TestRequired_MalformedHeader(t *testing.T) {
	r := newRouter("http://unused")
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Token abc") // not "Bearer ..."

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid authorization format")
}

func TestRequired_ValidToken_PropagatesIdentity(t *testing.T) {
	auth := stubAuth(`{"valid":true,"user_id":"u-123","role":"admin","name":"Alice"}`)
	defer auth.Close()

	r := newRouter(auth.URL)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer good-token")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "u-123")
	assert.Contains(t, w.Body.String(), "admin")
}

func TestRequired_InvalidToken(t *testing.T) {
	auth := stubAuth(`{"valid":false}`)
	defer auth.Close()

	r := newRouter(auth.URL)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer bad-token")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid or expired token")
}

func TestRequired_AuthServiceDown(t *testing.T) {
	// Closed server → validate() gets a connection error → 401.
	auth := stubAuth(`{"valid":true}`)
	url := auth.URL
	auth.Close()

	r := newRouter(url)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer any")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestOptional_NoHeaderPassesThrough(t *testing.T) {
	m := NewAuthMiddleware("http://unused", nil)
	r := gin.New()
	r.GET("/open", m.Optional(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/open", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
