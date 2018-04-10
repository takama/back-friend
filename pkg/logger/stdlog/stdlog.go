package stdlog

import (
	"log"
	"os"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/logger"
)

// UTC contains default UTC suffix
const UTC = "+0000 "

// New returns logger that is compatible with the Logger interface
func New(cfg *logger.Config) *Logger {
	var flags int
	prefix := "[" + config.ServiceName + ":" + cfg.Level.String() + "] "
	if cfg.Out == nil {
		cfg.Out = os.Stdout
	}
	if cfg.Err == nil {
		cfg.Err = os.Stderr
	}
	if cfg.Time {
		flags = log.Ldate | log.Ltime | log.Lmicroseconds
		if cfg.UTC {
			flags = flags | log.LUTC
		}
	}
	return &Logger{
		Level:  cfg.Level,
		Time:   cfg.Time,
		UTC:    cfg.UTC,
		stdlog: log.New(cfg.Out, prefix, flags),
		errlog: log.New(cfg.Err, prefix, flags),
	}
}

// Logger implements the Logger interface
// except of using logger.Fields
type Logger struct {
	logger.Level
	Time   bool
	UTC    bool
	stdlog *log.Logger
	errlog *log.Logger
}

// Debug logs a debug message
func (l *Logger) Debug(v ...interface{}) {
	if l.Level == logger.LevelDebug {
		l.setStdPrefix(logger.LevelDebug)
		l.printStd(v...)
	}
}

// Debugf logs a debug message with format
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.Level == logger.LevelDebug {
		l.setStdPrefix(logger.LevelDebug)
		l.printfStd(format, v...)
	}
}

// Info logs a info message
func (l *Logger) Info(v ...interface{}) {
	if l.Level <= logger.LevelInfo {
		l.setStdPrefix(logger.LevelInfo)
		l.printStd(v...)
	}
}

// Infof logs a info message with format
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.Level <= logger.LevelInfo {
		l.setStdPrefix(logger.LevelInfo)
		l.printfStd(format, v...)
	}
}

// Warn logs a warning message.
func (l *Logger) Warn(v ...interface{}) {
	if l.Level <= logger.LevelWarn {
		l.setStdPrefix(logger.LevelWarn)
		l.printStd(v...)
	}
}

// Warnf logs a warning message with format.
func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.Level <= logger.LevelWarn {
		l.setStdPrefix(logger.LevelWarn)
		l.printfStd(format, v...)
	}
}

// Error logs an error message
func (l *Logger) Error(v ...interface{}) {
	if l.Level <= logger.LevelError {
		l.setErrPrefix(logger.LevelError)
		l.printErr(v...)
	}
}

// Errorf logs an error message with format
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.Level <= logger.LevelError {
		l.setErrPrefix(logger.LevelError)
		l.printfErr(format, v...)
	}
}

// Fatal logs an error message followed by a call to os.Exit(1)
func (l *Logger) Fatal(v ...interface{}) {
	if l.Level <= logger.LevelFatal {
		l.setErrPrefix(logger.LevelFatal)
		l.printErr(v...)
	}
}

// Fatalf logs an error message with format followed by a call to ox.Exit(1)
func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l.Level <= logger.LevelFatal {
		l.setErrPrefix(logger.LevelFatal)
		l.printfErr(format, v...)
	}
}

func (l *Logger) printStd(v ...interface{}) {
	if l.Time && l.UTC {
		l.stdlog.Print(append([]interface{}{UTC}, v...)...)
	} else {
		l.stdlog.Print(v...)
	}
}

func (l *Logger) printfStd(format string, v ...interface{}) {
	if l.Time && l.UTC {
		l.stdlog.Printf("%s"+format, append([]interface{}{UTC}, v...)...)
	} else {
		l.stdlog.Printf(format, v...)
	}
}

func (l *Logger) printErr(v ...interface{}) {
	if l.Time && l.UTC {
		l.errlog.Print(append([]interface{}{UTC}, v...)...)
	} else {
		l.errlog.Print(v...)
	}
}

func (l *Logger) printfErr(format string, v ...interface{}) {
	if l.Time && l.UTC {
		l.errlog.Printf("%s"+format, append([]interface{}{UTC}, v...)...)
	} else {
		l.errlog.Printf(format, v...)
	}
}

func (l *Logger) setStdPrefix(level logger.Level) {
	l.stdlog.SetPrefix("[" + config.ServiceName + ":" + level.String() + "] ")
}

func (l *Logger) setErrPrefix(level logger.Level) {
	l.errlog.SetPrefix("[" + config.ServiceName + ":" + level.String() + "] ")
}
