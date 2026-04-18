package repository

import (
	"context"
	"github.com/mmtaee/ocserv-dashboard/api/pkg/request"
	"github.com/mmtaee/ocserv-dashboard/common/models"
	"github.com/mmtaee/ocserv-dashboard/common/ocserv/user"
	"github.com/mmtaee/ocserv-dashboard/common/pkg/database"
	"gorm.io/gorm"
	"time"
)

type ReportRepository struct {
	db                   *gorm.DB
	commonOcservUserRepo user.OcservUserInterface
}

type ReportRepositoryInterface interface {
	SessionLogs(ctx context.Context, pagination *request.Pagination, dateStart, dateEnd *time.Time) (*[]models.OcservUserSessionLog, int64, error)
}

func NewtReportRepository() *ReportRepository {
	return &ReportRepository{
		db:                   database.GetConnection(),
		commonOcservUserRepo: user.NewOcservUser(),
	}
}

func (r *ReportRepository) SessionLogs(
	ctx context.Context,
	pagination *request.Pagination,
	dateStart, dateEnd *time.Time,
) (*[]models.OcservUserSessionLog, int64, error) {
	var totalRecords int64

	query := r.db.WithContext(ctx).Model(&models.OcservUserSessionLog{})

	if dateStart != nil {
		query = query.Where("created_at >= ?", *dateStart)
	}

	if dateEnd != nil {
		query = query.Where("created_at < ?", dateEnd.AddDate(0, 0, 1))
	}

	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	var logs []models.OcservUserSessionLog
	if err := request.Paginator(ctx, query, pagination).
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return &logs, totalRecords, nil
}
