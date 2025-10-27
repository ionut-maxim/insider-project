package worker

import "errors"

var (
	ErrAlreadyRunning        = errors.New("worker already running")
	ErrAlreadyStopped        = errors.New("worker already stopped")
	ErrUnableToFetchMessages = errors.New("unable to fetch messages")
	ErrUnableToNotify        = errors.New("unable to notify")
	ErrUnableToUpdateStatus  = errors.New("unable to update status")
	ErrUnableToUpdateCache   = errors.New("unable to update cache")
)
