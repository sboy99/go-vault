package domain

import "errors"

var (
	ErrBackupNotFound         = errors.New("backup not found")
	ErrJobNotFound            = errors.New("job not found")
	ErrRestoreConfirmMismatch = errors.New("confirm must match the configured database name")
	ErrJobAlreadyRunning      = errors.New("another job is already running")
	ErrBackupNotRestorable    = errors.New("backup is not in a restorable state")
)
