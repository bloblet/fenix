package utils

import (
	"fmt"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

var logger *log.Logger
var once = sync.Once{}

func Log() *log.Logger {
	once.Do(setupLogger)
	return logger
}

type FenixFormatter struct {
}

func (f *FenixFormatter) Format(entry *log.Entry) ([]byte, error) {
	var level string
	var keyColor *color.Color

	switch entry.Level {
	case log.TraceLevel:
		level = "TRAC"
		keyColor = color.New(color.FgHiBlack, color.BgHiCyan)
	case log.DebugLevel:
		level = "DEBU"
		keyColor = color.New(color.FgHiBlack, color.BgHiMagenta)
	case log.InfoLevel:
		level = "INFO"
		keyColor = color.New(color.FgHiBlack, color.BgHiGreen)
	case log.WarnLevel:
		level = "WARN"
		keyColor = color.New(color.FgHiBlack, color.BgHiYellow)
	case log.ErrorLevel:
		level = "ERRO"
		keyColor = color.New(color.FgHiBlack, color.BgRed)
	case log.FatalLevel:
		level = "FATA"
		keyColor = color.New(color.FgHiBlack, color.BgHiRed)
	case log.PanicLevel:
		level = "PANI"
		keyColor = color.New(color.FgHiBlack, color.BgHiWhite)
	}

	function := strings.Split(entry.Caller.Function, "/")
	function_name := function[len(function)-1]

	d := fmt.Sprintf(
		`[%v] %v %v      %v
    %v %v
`,
		entry.Time.Format(time.RFC3339), color.HiMagentaString("[%v]", level), color.HiGreenString("%20v", function_name), color.GreenString("%v:%v", entry.Caller.File, entry.Caller.Line),
		keyColor.Sprint("[msg]"),
		color.MagentaString(entry.Message),
	)

	for key, value := range entry.Data {
		if value == nil || value == "" {
			d = fmt.Sprintf("%v    %v %v\n", d, keyColor.Sprintf("[%v]", key), color.MagentaString("nil"))
		} else {
			d = fmt.Sprintf("%v    %v %v\n", d, keyColor.Sprintf("[%v]", key), color.MagentaString("%v", value))
		}
	}
	d = fmt.Sprintf("%v\n", d)

	return []byte(d), nil
}
func TestLogger() {
	Log().WithFields(
		log.Fields{
			"field":     "field1",
			"nilField":  "",
			"nilField2": nil,
		},
	).Trace("Testing loggs")
	Log().WithFields(
		log.Fields{
			"field":     "field1",
			"nilField":  "",
			"nilField2": nil,
		},
	).Debug("Testing loggs")
	Log().WithFields(
		log.Fields{
			"field":     "field1",
			"nilField":  "",
			"nilField2": nil,
		},
	).Info("Testing loggs")
	Log().WithFields(
		log.Fields{
			"field":     "field1",
			"nilField":  "",
			"nilField2": nil,
		},
	).Warn("Testing loggs")
	Log().WithFields(
		log.Fields{
			"field":     "field1",
			"nilField":  "",
			"nilField2": nil,
		},
	).Error("Testing loggs")

}

func setupLogger() {
	l := log.New()
	l.ReportCaller = true
	var config = LoadConfig("fenix.yml")
	l.SetFormatter(&FenixFormatter{})

	level, err := log.ParseLevel(config.Logger.LogLevel)
	if err != nil {
		l.WithFields(log.Fields{
			"level": level,
			"err":   err,
		}).Panic("Failed to parse log.LogLevel")
	}
	l.SetLevel(level)

	logger = l
}
