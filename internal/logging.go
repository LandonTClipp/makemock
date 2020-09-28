package internal

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	LogKeyCommand = "command"
	LogKeyVersion = "version"
)

// GetNewLogger returns a newly allocated log object
func GetNewLogger(disableColor bool, level string) (zerolog.Logger, error) {
	var isTerminal bool
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		isTerminal = true
	}
	var color bool
	if disableColor {
		color = false
	} else {
		color = isTerminal
	}

	writer := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    !color,
		TimeFormat: time.RFC3339,
	}

	log := zerolog.New(writer)
	log = log.Hook(timeHook{})
	log = log.With().Str(LogKeyVersion, Version).Logger()

	parsedLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return log, errors.Wrapf(err, "failed to parse level")
	}
	zerolog.SetGlobalLevel(parsedLevel)

	return log, nil
}

type timeHook struct{}

func (t timeHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	e.Time("time", time.Now())
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// StackAndFail prints the stack trace from err (if any) then exits the program
func StackAndFail(err error) {
	cause := errors.Cause(err)
	fmt.Printf("%v", cause)
	stacker, ok := err.(stackTracer)
	if ok {
		fmt.Printf("%+v\n", stacker.StackTrace())
	}
	os.Exit(1)
}
