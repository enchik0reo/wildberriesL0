package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

type Lgr struct {
	*logrus.Logger
}

func New() *Lgr {
	l := logrus.New()
	l.SetOutput(io.Discard)
	l.SetReportCaller(true)

	if err := os.MkdirAll("logs", 0744); err != nil {
		panic(err)
	}

	fl, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	l.AddHook(&Hook1{
		Writer:    []io.Writer{fl},
		LogLevels: logrus.AllLevels,
	})

	l.AddHook(&Hook2{
		Writer:    []io.Writer{os.Stdout},
		LogLevels: []logrus.Level{logrus.InfoLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel},
	})

	return &Lgr{
		Logger: l,
	}
}

type Hook1 struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (h *Hook1) Levels() []logrus.Level {
	return h.LogLevels
}

func (h *Hook1) Fire(entry *logrus.Entry) error {
	entry.Logger.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			fName := path.Base(f.File)
			function = fmt.Sprintf("%s()", f.Function)
			file = fmt.Sprintf("%s:%d", fName, f.Line)
			return function, file
		},
		FullTimestamp: true,
	}

	str, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range h.Writer {
		_, err = w.Write([]byte(str))
	}
	return err
}

type Hook2 struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (h *Hook2) Levels() []logrus.Level {
	return h.LogLevels
}

func (h *Hook2) Fire(entry *logrus.Entry) error {
	entry.Logger.Formatter = &logrus.TextFormatter{
		ForceColors: true,
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			fName := path.Base(f.File)
			function = fmt.Sprintf("%s()\nMessage:", f.Function)
			file = fmt.Sprintf(" | %s:%d |", fName, f.Line)
			return function, file
		},
		FullTimestamp: true,
	}

	str, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range h.Writer {
		_, err = w.Write([]byte(str + "\n"))
	}
	return err
}
