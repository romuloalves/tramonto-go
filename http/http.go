package http

import (
	"io/ioutil"
	"net/http"
	"sync"

	"gitlab.com/tramonto-one/go-tramonto/entities"

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
	// Sets server in RELEASE mode
	gin.SetMode(gin.ReleaseMode)

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

	h.server.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	if err := h.server.Run(httpPort); err != nil {
		return err
	}

	return nil
}

// AddPostArtifact registers and calls the function to add a new artifacts to a test
func (h *OneHTTP) AddPostArtifact(callback func(ipns, name, description string, file []byte, headers map[string][]string) ([]byte, error)) {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.server.POST("/artifacts/:ipns", func(c *gin.Context) {
		ipns := c.Param("ipns")

		form, err := c.MultipartForm()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		name := form.Value["name"][0]
		description := form.Value["description"][0]
		file, err := c.FormFile("artifact")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		headers := (map[string][]string)(file.Header)

		fileReader, err := file.Open()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		content, err := ioutil.ReadAll(fileReader)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		response, err := callback(ipns, name, description, content, headers)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.String(http.StatusOK, string(response))
	})
}

// AddGetArtifact registers and calls the function to download an artifact from a test
func (h *OneHTTP) AddGetArtifact(callback func(ipns, artifactHash string) (entities.Artifact, []byte, error)) {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.server.GET("/artifacts/:ipns/:artifactHash", func(c *gin.Context) {
		ipns := c.Param("ipns")
		artifactHash := c.Param("artifactHash")

		artifact, content, err := callback(ipns, artifactHash)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		for key, value := range artifact.Headers {
			c.Header(key, value[0])
		}

		contentType := artifact.Headers["Content-Type"][0]

		c.Data(http.StatusOK, contentType, content)
	})
}
