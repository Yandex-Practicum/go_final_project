package handlers

import (
	"errors"
	"log/slog"

	"github.com/wissio/go_final_project/internal/storage/sqlite"
)

type Handlers struct {
	s   *sqlite.Storage
	log *slog.Logger
}

const (
	DateLayout    = "20060102"
	DateDotLayout = "02.01.2006"
	taskLimit     = 10
)

var ErrTaskNotFound = errors.New("задача не найдена")
