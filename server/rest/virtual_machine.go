package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/arxeiss/sample-terraform-provider/entities"
)

func (h *HTTPServer) getVM(c *gin.Context) {
	id := getPathID(c)
	if id == 0 {
		return // error was sent already
	}

	vm, err := h.db.VirtualMachines.FindByID(id)
	if err != nil {
		h.serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, vm)
}

func (h *HTTPServer) createVM(c *gin.Context) {
	vm := &entities.VirtualMachine{}
	if err := c.ShouldBindJSON(&vm); err != nil {
		h.bindingError(c, err)
		return
	}
	if err := vm.Validate(); err != nil {
		h.validationError(c, err)
		return
	}

	vm, err := h.db.VirtualMachines.Save(vm)
	if err != nil {
		h.serverError(c, err)
		return
	}
	c.JSON(http.StatusCreated, vm)
}

func (h *HTTPServer) editVM(c *gin.Context) {
	vm := &entities.VirtualMachine{}
	if err := c.ShouldBindJSON(&vm); err != nil {
		h.bindingError(c, err)
		return
	}
	if err := vm.Validate(); err != nil {
		h.validationError(c, err)
		return
	}

	vm.ID = getPathID(c)
	if vm.ID == 0 {
		return // error was sent already
	}

	vm, err := h.db.VirtualMachines.Save(vm)
	if err != nil {
		h.serverError(c, err)
		return
	}
	c.JSON(http.StatusCreated, vm)
}

func (h *HTTPServer) deleteVM(c *gin.Context) {
	id := getPathID(c)
	if id == 0 {
		return // error was sent already
	}

	err := h.db.VirtualMachines.Delete(id)
	if err != nil {
		h.serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}
