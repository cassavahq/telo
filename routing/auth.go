package routing
import (
	"net/http"
	"github.com/labstack/echo"
)

func AuthLogin(c echo.Context) error {
	username := c.FormValue("username")

	return c.String(http.StatusOK, username)
}