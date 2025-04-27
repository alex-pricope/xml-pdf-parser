package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func BoostrapLogger() {
	Log = &logrus.Logger{
		Out:   nil,
		Hooks: nil,
		Formatter: &logrus.TextFormatter{
			DisableColors:    false,
			DisableQuote:     false,
			DisableTimestamp: false,
			FullTimestamp:    true,
			CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
				filename := filepath.Base(f.File)
				return f.Function, fmt.Sprintf("%s:%d", filename, f.Line)
			},
		},
		ReportCaller: true,
		Level:        logrus.DebugLevel,
		ExitFunc:     nil,
	}

	Log.Out = os.Stdout
}
