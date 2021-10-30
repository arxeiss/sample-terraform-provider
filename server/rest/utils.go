package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/arxeiss/sample-terraform-provider/server/database"
)

func (*HTTPServer) bindingError(c *gin.Context, err error) {
	je := &json.UnmarshalTypeError{}
	if errors.As(err, &je) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "cannot unmarshal " + je.Value + " into field " + je.Field + " of type " + je.Type.String(),
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
}

func (*HTTPServer) validationError(c *gin.Context, err error) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"status": http.StatusUnprocessableEntity,
		"error":  err.Error(),
	})
}

func (h *HTTPServer) serverError(c *gin.Context, err error) {
	if database.IsNotFoundError(err) {
		h.error404(c)
		return
	}
	if database.IsUniqueConstraintError(err) {
		h.validationError(c, errors.New("given name already exists"))
		return
	}
	h.log.Error(err)
	c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError})
}

func (*HTTPServer) error404(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "error": "Not found"})
}

func authorization(token string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		auth := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")
		if auth != token {
			// Credentials doesn't match, we return 401 and abort handlers chain.
			ctx.Header("WWW-Authenticate", "Basic realm=Authorization required")
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"status": http.StatusUnauthorized, "error": "Authorization required"},
			)
			return
		}
	}
}

func (h *HTTPServer) loggerMiddleware(c *gin.Context) {
	// Start timer
	start := time.Now()

	// Process request
	c.Next()

	status := c.Writer.Status()
	if status == http.StatusInternalServerError {
		// Don't log twice server error
		return
	}
	path := c.Request.URL.Path
	if raw := c.Request.URL.RawQuery; raw != "" {
		path = path + "?" + raw
	}

	h.log.
		WithField("method", c.Request.Method).
		WithField("latency", time.Since(start).String()).
		WithField("status", status).
		Info(path)
}

func getPathID(c *gin.Context) int64 {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "cannot parse id in URL"})
	}
	return id
}
