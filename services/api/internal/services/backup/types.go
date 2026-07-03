package backup

import (
	"time"

	"github.com/mmtaee/ocserv-dashboard/common/models"
)

const (
	BackupVersion = "2.0.0"
	BackupTypeUsers = "users"
	BackupTypeGroups = "groups"
)

type BackupMetadata struct {
	Version          string    `json:"version"`
	DashboardVersion string    `json:"dashboard_version"`
	BackupType       string    `json:"backup_type"`
	Hostname         string    `json:"hostname"`
	CreatedAt        time.Time `json:"created_at"`
	Compression      string    `json:"compression"`
	Checksum         string    `json:"checksum,omitempty"`
}

type UserBackup struct {
	Metadata BackupMetadata      `json:"metadata"`
	Users    []models.OcservUser `json:"users"`
}

type GroupBackup struct {
	Metadata     BackupMetadata            `json:"metadata"`
	DefaultGroup *models.OcservGroupConfig `json:"default_group"`
	Groups       []models.OcservGroup      `json:"groups"`
}
