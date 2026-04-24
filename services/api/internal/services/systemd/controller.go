package systemd

import (
	"github.com/labstack/echo/v4"
	"github.com/mmtaee/ocserv-dashboard/api/internal/repository"
	"github.com/mmtaee/ocserv-dashboard/api/pkg/request"
	"net/http"
)

type Controller struct {
	request request.CustomRequestInterface
	systemd repository.SystemdRepositoryInterface
}

func New() *Controller {
	return &Controller{
		request: request.NewCustomRequest(),
		systemd: repository.NewSystemdRepository(),
	}
}

// Status
// @Summary      Ocserv systemctl status
// @Description  Ocserv systemctl status
// @Tags         Systemd
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Failure      429 {object} middlewares.TooManyRequests
// @Success      200 {object}  OcservSystemdStatus
// @Router       /systemd/status [get]
func (ctl *Controller) Status(c echo.Context) error {
	//if os.Getenv("SYSTEMD") != "true" {
	//	return ctl.request.BadRequest(c, errors.New("systemd is not running"))
	//}

	statusLog, err := ctl.systemd.Status(c.Request().Context())
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}

	output := ParseSystemctlShow(statusLog)
	return c.JSON(http.StatusOK, output)
}
