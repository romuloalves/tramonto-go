package http

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

const httpPort = ":3000"

// OneHTTP represents a HTTP server to Tramonto One
type OneHTTP struct {
	server *gin.Engine
	mux    *sync.Mutex
}

// InitializeHTTPServer initializes the new HTTP server
func InitializeHTTPServer() (*OneHTTP, error) {
	r := gin.Default()
	http := &OneHTTP{
		server: r,
		mux:    new(sync.Mutex),
	}

	return http, nil
}

// Start starts the server
func (h *OneHTTP) Start() error {
	h.mux.Lock()
	defer h.mux.Unlock()

	if err := h.server.Run(httpPort); err != nil {
		return err
	}

	return nil
}

// AddPostArtifact registers and calls the function to add a new artifacts to a test
func (h *OneHTTP) AddPostArtifact(callback func(ipns, name, description string) ([]byte, error)) {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.server.POST("/:ipns", func(c *gin.Context) {
		ipns := c.Param("ipns")

		form, err := c.MultipartForm()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		name := form.Value["name"][0]
		description := form.Value["description"][0]
		// file := form.File["artifact[]"][0]

		response, err := callback(ipns, name, description)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.String(http.StatusOK, string(response))
	})
}

// AddGetArtifact registers and calls the function to download an artifact from a test
func (h *OneHTTP) AddGetArtifact(callback func(ipns, artifactHash string) ([]byte, error)) {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.server.GET("/:ipns/:artifactHash", func(c *gin.Context) {
		ipns := c.Param("ipns")
		artifactHash := c.Param("artifactHash")

		response, err := callback(ipns, artifactHash)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.String(http.StatusOK, string(response))
	})
}
