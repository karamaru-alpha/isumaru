package mysql

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type Handler interface {
	TailSlowQueryLog(c echo.Context) error
}

type handler struct {
	slowQueryLogPath string
}

func New(slowQueryLogPath string) Handler {
	return &handler{
		slowQueryLogPath,
	}
}

func (h *handler) TailSlowQueryLog(c echo.Context) error {
	seconds, err := strconv.Atoi(c.QueryParam("seconds"))
	if err != nil {
		return err
	}

	file, err := os.Open(h.slowQueryLogPath)
	if err != nil {
		return fmt.Errorf("failed to open slow query log file. err=%+v", err)
	}
	defer file.Close()

	start, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("failed to seek slow query log file at start. err=%+v", err)
	}

	time.Sleep(time.Second * time.Duration(seconds))

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat slow query log file. err=%+v", err)
	}

	reader := io.LimitReader(file, stat.Size()-start)
	if err := c.Stream(http.StatusOK, "application/json", reader); err != nil {
		return fmt.Errorf("failed to stream slow query log. err=%+v", err)
	}

	return nil
}
