package statistics

import (
	"github.com/labstack/echo/v4"
	"github.com/mmtaee/ocserv-dashboard/api/pkg/routing/middlewares"
)

func Routes(e *echo.Group) {
	ctl := New()
	g := e.Group("/reports", middlewares.AuthMiddleware())

	g.GET("/session_logs", ctl.SessionLogs)
}
