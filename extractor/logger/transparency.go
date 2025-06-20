package logger

import (
	"github.com/warnakulasuriya-fds-e23/go-sourceafis-fork/features"
)

type TransparencyLogger interface {
	Log(key string, data interface{}) error
	LogSkeleton(keyword string, skeleton *features.Skeleton) error
}
