package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/arxeiss/sample-terraform-provider/entities"
)

func (h *HTTPServer) getStorage(c *gin.Context) {
	id := getPathID(c)
	if id == 0 {
		return // error was sent already
	}

	s, err := h.db.Storages.FindByID(id)
	if err != nil {
		h.serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *HTTPServer) createStorage(c *gin.Context) {
	s := &entities.Storage{}
	if err := c.ShouldBindJSON(&s); err != nil {
		h.bindingError(c, err)
		return
	}
	if err := s.Validate(); err != nil {
		h.validationError(c, err)
		return
	}

	s, err := h.db.Storages.Save(s) //nolint:ifshort
	if err != nil {
		h.serverError(c, err)
		return
	}
	c.JSON(http.StatusCreated, s)
}

func (h *HTTPServer) editStorage(c *gin.Context) {
	s := &entities.Storage{}
	if err := c.ShouldBindJSON(&s); err != nil {
		h.bindingError(c, err)
		return
	}
	if err := s.Validate(); err != nil {
		h.validationError(c, err)
		return
	}

	s.ID = getPathID(c)
	if s.ID == 0 {
		return // error was sent already
	}

	s, err := h.db.Storages.Save(s) //nolint:ifshort
	if err != nil {
		h.serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *HTTPServer) deleteStorage(c *gin.Context) {
	id := getPathID(c)
	if id == 0 {
		return // error was sent already
	}

	err := h.db.Storages.Delete(id)
	if err != nil {
		h.serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}
