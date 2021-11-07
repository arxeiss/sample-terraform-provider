package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/arxeiss/sample-terraform-provider/entities"
)

func (h *HTTPServer) getNetwork(c *gin.Context) {
	id := getPathID(c)
	if id == 0 {
		return // error was sent already
	}

	n, err := h.db.Networks.FindByID(id)
	if err != nil {
		h.serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, n)
}

func (h *HTTPServer) createNetwork(c *gin.Context) {
	n := &entities.Network{}
	if err := c.ShouldBindJSON(&n); err != nil {
		h.bindingError(c, err)
		return
	}
	if err := n.Validate(); err != nil {
		h.validationError(c, err)
		return
	}

	n, err := h.db.Networks.Save(n) //nolint:ifshort
	if err != nil {
		h.serverError(c, err)
		return
	}
	c.JSON(http.StatusCreated, n)
}

func (h *HTTPServer) editNetwork(c *gin.Context) {
	n := &entities.Network{}
	if err := c.ShouldBindJSON(&n); err != nil {
		h.bindingError(c, err)
		return
	}
	if err := n.Validate(); err != nil {
		h.validationError(c, err)
		return
	}

	n.ID = getPathID(c)
	if n.ID == 0 {
		return // error was sent already
	}

	n, err := h.db.Networks.Save(n) //nolint:ifshort
	if err != nil {
		h.serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, n)
}

func (h *HTTPServer) deleteNetwork(c *gin.Context) {
	id := getPathID(c)
	if id == 0 {
		return // error was sent already
	}

	err := h.db.Networks.Delete(id)
	if err != nil {
		h.serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}
