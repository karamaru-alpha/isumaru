package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type PProfHandler interface {
	GetPProf(c echo.Context) error
}

type pprofHandler struct {
	pprofDirFormat string
}

func NewPProfHandler(pprofDirFormat string) PProfHandler {
	return &pprofHandler{pprofDirFormat}
}

func (h *pprofHandler) GetPProf(c echo.Context) error {
	entryID := c.Param("entryID")
	targetID := c.Param("targetID")
	dir := fmt.Sprintf(h.pprofDirFormat, entryID)
	path := fmt.Sprintf("%s/%s.pprof", dir, targetID)
	c.Response().Header().Set(echo.HeaderContentType, "application/octet-stream")
	return c.File(path)
}
