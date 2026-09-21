package http

import (
	"errors"
	nethttp "net/http"

	"github.com/sboy99/go-vault/internal/domain"
)

func writeDomainError(w nethttp.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrBackupNotFound), errors.Is(err, domain.ErrJobNotFound):
		writeError(w, nethttp.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrRestoreConfirmMismatch), errors.Is(err, domain.ErrBackupNotRestorable):
		writeError(w, nethttp.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrJobAlreadyRunning):
		writeError(w, nethttp.StatusConflict, err.Error())
	default:
		writeError(w, nethttp.StatusInternalServerError, err.Error())
	}
}
