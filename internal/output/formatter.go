package output

import (
	"io"

	"github.com/joaovictornsv/cards-cli/internal/buildinfo"
	"github.com/joaovictornsv/cards-cli/internal/config"
)

type Formatter interface {
	PrintConfig(w io.Writer, cfg config.Config) error
	PrintVersion(w io.Writer, info buildinfo.Info) error
}

func New(json bool) Formatter {
	if json {
		return JSONFormatter{}
	}
	return TableFormatter{}
}
