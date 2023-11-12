package logger

import (
	"github.com/go-mods/zerolog-quick/console/colored"
	"github.com/go-mods/zerolog-quick/console/plain"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

var isTerm = isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())

func init() {
	log.Logger = NewConsoleLogger().Level(zerolog.InfoLevel)
}

func NewConsoleLogger() zerolog.Logger {
	// If the output is a terminal, colorize the output
	if isTerm {
		return colored.Message
	} else {
		return plain.Message
	}
}
