package statistics

import (
	"github.com/mmtaee/ocserv-dashboard/api/pkg/request"
	"github.com/mmtaee/ocserv-dashboard/common/models"
)

type SessionLogsData struct {
	DateStart string `json:"date_start" query:"date_start" validate:"omitempty" example:"2025-1-31"`
	DateEnd   string `json:"date_end" query:"date_end" validate:"omitempty" example:"2025-12-31"`
}

type SessionLogsResponse struct {
	Meta   request.Meta                   `json:"meta" validate:"required"`
	Result *[]models.OcservUserSessionLog `json:"result" validate:"omitempty"`
}
