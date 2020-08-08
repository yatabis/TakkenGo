package line

import (
	"net/http"

	"github.com/labstack/echo"
)

func Callback(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
