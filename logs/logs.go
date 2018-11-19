/**
 * Custom logger that wraps Beego logger, add some features as below:
 * 1. Send msg to sentry for error with level Critical or less
 */
package logs

import (
	"github.com/astaxie/beego/logs"
	"github.com/getsentry/raven-go"
	"strings"
	"fmt"
)

var (
	beeLogger    *logs.BeeLogger
	sentryClient *raven.Client
	serviceName   string
)

const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

// Name for adapter with beego official support
const (
	AdapterConsole   = "console"
	AdapterFile      = "file"
	AdapterMultiFile = "multifile"
	AdapterMail      = "smtp"
	AdapterConn      = "conn"
	AdapterEs        = "es"
	AdapterJianLiao  = "jianliao"
	AdapterSlack     = "slack"
	AdapterAliLS     = "alils"
)

type ErrWriter struct {

}

func (w ErrWriter)Write(p []byte) (n int, err error) {
	Emergency(string(p))
	return len(p), nil
}

func init() {
	beeLogger = logs.NewLogger()
}

func SetupSentry(sentryDsn string, sName string) {
	var err error
	if sentryClient, err = raven.New(sentryDsn); err != nil {
		panic(err)
		Critical("setup sentry fails: %s", err.Error())
	}
	serviceName = sName
}

func SetLevel(level int)  {
	beeLogger.SetLevel(level)
}

// SetLogger sets a new logger.
func SetLogger(adapter string, config ...string) error {
	return beeLogger.SetLogger(adapter, config...)
}

func GetBeeLogger() *logs.BeeLogger {
	return beeLogger
}

// Emergency logs a message at emergency level.
func Emergency(f interface{}, v ...interface{}) {
	var msg = formatLog(f, v...)
	beeLogger.Emergency(msg)

	if sentryClient != nil{
		sentryClient.CaptureMessage(msg, sentryTags())
	}
}

// Alert logs a message at alert level.
func Alert(f interface{}, v ...interface{}) {
	beeLogger.Alert(formatLog(f, v...))
}

// Critical logs a message at critical level.
func Critical(f interface{}, v ...interface{}) {
	var msg = formatLog(f, v...)
	beeLogger.Critical(msg)

	if sentryClient != nil{
		sentryClient.CaptureMessage(msg, sentryTags())
	}
}

// Error logs a message at error level.
func Error(f interface{}, v ...interface{}) {
	beeLogger.Error(formatLog(f, v...))
}

// Warning logs a message at warning level.
func Warning(f interface{}, v ...interface{}) {
	beeLogger.Warn(formatLog(f, v...))
}

// Warn compatibility alias for Warning()
func Warn(f interface{}, v ...interface{}) {
	beeLogger.Warn(formatLog(f, v...))
}

// Notice logs a message at notice level.
func Notice(f interface{}, v ...interface{}) {
	beeLogger.Notice(formatLog(f, v...))
}

// Informational logs a message at info level.
func Informational(f interface{}, v ...interface{}) {
	beeLogger.Info(formatLog(f, v...))
}

// Info compatibility alias for Warning()
func Info(f interface{}, v ...interface{}) {
	beeLogger.Info(formatLog(f, v...))
}

// Debug logs a message at debug level.
func Debug(f interface{}, v ...interface{}) {
	beeLogger.Debug(formatLog(f, v...))
}

// Trace logs a message at trace level.
// compatibility alias for Warning()
func Trace(f interface{}, v ...interface{}) {
	beeLogger.Trace(formatLog(f, v...))
}

func formatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return fmt.Sprintf(msg, v...)
}

func sentryTags() map[string]string {
	var tags = make(map[string]string)

	if serviceName != "" {
		tags["service_name"] = serviceName
	}

	return tags
}
