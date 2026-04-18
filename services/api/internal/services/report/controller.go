package statistics

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mmtaee/ocserv-dashboard/api/internal/repository"
	"github.com/mmtaee/ocserv-dashboard/api/pkg/request"
	"net/http"
	"time"
)

type Controller struct {
	request    request.CustomRequestInterface
	reportRepo repository.ReportRepositoryInterface
}

func New() *Controller {
	return &Controller{
		request:    request.NewCustomRequest(),
		reportRepo: repository.NewtReportRepository(),
	}
}

// SessionLogs 	 Ocserv session logs
//
// @Summary      Ocserv session logs
// @Description  Ocserv session logs
// @Tags         Report
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 page query int false "Page number, starting from 1" minimum(1)
// @Param 		 size query int false "Number of items per page" minimum(1) maximum(100) name(size)
// @Param 		 order query string false "Field to order by"
// @Param 		 sort query string false "Sort order, either ASC or DESC" Enums(ASC, DESC)
// @Param 		 date_start query string false "date_start"
// @Param 		 date_end query string false "date_end"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      200  {object} SessionLogsResponse
// @Router       /reports/session_logs [get]
func (ctl *Controller) SessionLogs(c echo.Context) error {
	var data SessionLogsData
	if err := c.Bind(&data); err != nil {
		return ctl.request.BadRequest(c, err)
	}

	pagination := ctl.request.Pagination(c)

	var startDate, endDate *time.Time

	if data.DateStart != "" {
		t, err := time.Parse("2006-01-02", data.DateStart)
		if err != nil {
			return ctl.request.BadRequest(c, fmt.Errorf("invalid date_start: %w", err))
		}
		startDate = &t
	}

	if data.DateEnd != "" {
		t, err := time.Parse("2006-01-02", data.DateEnd)
		if err != nil {
			return ctl.request.BadRequest(c, fmt.Errorf("invalid date_end: %w", err))
		}
		t = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		endDate = &t
	}

	logs, total, err := ctl.reportRepo.SessionLogs(c.Request().Context(), pagination, startDate, endDate)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}

	return c.JSON(http.StatusOK, SessionLogsResponse{
		Meta: request.Meta{
			Page:         pagination.Page,
			TotalRecords: total,
			PageSize:     pagination.PageSize,
		},
		Result: logs,
	})
}
