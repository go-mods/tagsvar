// Package logger configures the zerolog logger of tagsvar.
//
// Messages are written to stderr, stdout being reserved for the output of the commands.
// They are colorized when the output is a terminal.
package logger

import (
	"io"
	"os"

	"github.com/go-mods/zerolog-quick/console/colored"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = New(os.Stderr).Level(zerolog.InfoLevel)
}

// New returns a console logger which writes only the messages to w.
// Messages are colorized when w is a terminal.
func New(w io.Writer) zerolog.Logger {
	output := zerolog.ConsoleWriter{
		Out:        w,
		NoColor:    true,
		PartsOrder: []string{zerolog.MessageFieldName},
	}
	if isTerminal(w) {
		output.FormatExtra = colored.Colorize
	}
	return zerolog.New(output)
}

// isTerminal reports whether w is a terminal
func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
}
