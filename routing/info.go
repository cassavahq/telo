package routing
import (
	"net/http"
	"github.com/labstack/echo"
)

// Info API
func Info(c echo.Context) error {
	return c.String(http.StatusOK, "INFO")
}