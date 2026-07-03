package backup

import "errors"

var (

	ErrInvalidBackup = errors.New("invalid backup file")

	ErrInvalidJSON = errors.New("invalid json")

	ErrInvalidChecksum = errors.New("invalid checksum")

	ErrInvalidVersion = errors.New("unsupported backup version")

	ErrInvalidCompression = errors.New("invalid compression")

	ErrBackupCorrupted = errors.New("backup is corrupted")
)
