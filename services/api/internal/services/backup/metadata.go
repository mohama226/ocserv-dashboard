package backup

import (
	"os"
	"time"
)

func NewMetadata(t string) BackupMetadata {

	host, _ := os.Hostname()

	return BackupMetadata{

		Version: BackupVersion,

		DashboardVersion: BackupVersion,

		BackupType: t,

		Hostname: host,

		CreatedAt: time.Now().UTC(),

		Compression: "gzip",
	}
}
