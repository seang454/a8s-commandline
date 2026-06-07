package runtime

import (
	"io"

	"github.com/yourname/a8s/internal/api"
	internalauth "github.com/yourname/a8s/internal/auth"
	"github.com/yourname/a8s/internal/config"
	"github.com/yourname/a8s/internal/output"
)

type Runtime struct {
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer

	Config  config.Resolved
	API     *api.Client
	Auth    *internalauth.Manager
	Printer output.Printer
}
