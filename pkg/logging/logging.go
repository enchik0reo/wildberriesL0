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
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			fName := path.Base(f.File)
			function = fmt.Sprintf("%s()", f.Function)
			file = fmt.Sprintf("%s:%d", fName, f.Line)
			return function, file
		},
		FullTimestamp: true,
	}

	if err := os.MkdirAll("logs", 0744); err != nil {
		panic(err)
	}

	fl, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	l.AddHook(&Hooker1{
		Writer:    []io.Writer{fl},
		LogLevels: logrus.AllLevels,
	})

	l.AddHook(&Hooker2{
		Writer:    []io.Writer{os.Stdout},
		LogLevels: []logrus.Level{logrus.InfoLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel},
	})

	return &Lgr{
		Logger: l,
	}
}

type Hooker1 struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (h *Hooker1) Levels() []logrus.Level {
	return h.LogLevels
}

func (h *Hooker1) Fire(entry *logrus.Entry) error {
	str, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range h.Writer {
		_, err = w.Write([]byte(str))
	}
	return err
}

type Hooker2 struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (h *Hooker2) Levels() []logrus.Level {
	return h.LogLevels
}

func (h *Hooker2) Fire(entry *logrus.Entry) error {
	str, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range h.Writer {
		_, err = w.Write([]byte(str + "\n"))
	}
	return err
}
