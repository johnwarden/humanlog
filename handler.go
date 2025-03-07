package humanlog

import (
	"github.com/humanlogio/humanlog/internal/pkg/config"
	"github.com/kr/logfmt"
)

// Handler can recognize it's log lines, parse them and prettify them.
type Handler interface {
	CanHandle(line []byte) bool
	Prettify(skipUnchanged bool) []byte
	logfmt.Handler
}

var DefaultOptions = &HandlerOptions{
	TimeFields:    []string{"time", "ts", "@timestamp", "timestamp"},
	MessageFields: []string{"message", "msg"},
	LevelFields:   []string{"level", "lvl", "loglevel", "severity"},
}

type HandlerOptions struct {
	TimeFields    []string
	MessageFields []string
	LevelFields   []string
}

var _ = HandlerOptionsFrom(config.DefaultConfig) // ensure it's valid

func HandlerOptionsFrom(cfg config.Config) *HandlerOptions {
	opts := DefaultOptions
	if cfg.TimeFields != nil {
		opts.TimeFields = *cfg.TimeFields
	}
	if cfg.MessageFields != nil {
		opts.MessageFields = *cfg.MessageFields
	}
	if cfg.LevelFields != nil {
		opts.LevelFields = *cfg.LevelFields
	}
	return opts
}
