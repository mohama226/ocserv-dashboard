package backup

import (
	"os"
	"path/filepath"
)

const BackupDirectory = "/opt/ocserv-dashboard/backups"

func EnsureBackupDirectory() error {

	return os.MkdirAll(BackupDirectory, 0755)

}

func BackupFile(name string) string {

	return filepath.Join(BackupDirectory, name)

}
