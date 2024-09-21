package logging

import (
	"io"
	"log"
	"os"
	"sync"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

type logger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
	level       LogLevel
}

var (
	singleton *logger
	once      sync.Once
)

func GetLogger() *logger {
	once.Do(func() {
		singleton = newLogger(os.Stdout, INFO)
	})
	return singleton
}

func newLogger(out io.Writer, level LogLevel) *logger {
	return &logger{
		debugLogger: log.New(out, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		infoLogger:  log.New(out, "INFO:  ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(out, "WARN:  ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(out, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		fatalLogger: log.New(out, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
		level:       level,
	}
}

func SetLogLevel(level LogLevel) {
	GetLogger().level = level
}

func SetOutput(out io.Writer) {
	logger := GetLogger()
	logger.debugLogger.SetOutput(out)
	logger.infoLogger.SetOutput(out)
	logger.warnLogger.SetOutput(out)
	logger.errorLogger.SetOutput(out)
	logger.fatalLogger.SetOutput(out)
}

func (l *logger) Debug(v ...interface{}) {
	if l.level <= DEBUG {
		l.debugLogger.Println(v...)
	}
}

func (l *logger) Debugf(format string, v ...interface{}) {
	if l.level <= DEBUG {
		l.debugLogger.Printf(format, v...)
	}
}

func (l *logger) Info(v ...interface{}) {
	if l.level <= INFO {
		l.infoLogger.Println(v...)
	}
}

func (l *logger) Infof(format string, v ...interface{}) {
	if l.level <= INFO {
		l.infoLogger.Printf(format, v...)
	}
}

func (l *logger) Warn(v ...interface{}) {
	if l.level <= WARN {
		l.warnLogger.Println(v...)
	}
}

func (l *logger) Warnf(format string, v ...interface{}) {
	if l.level <= WARN {
		l.warnLogger.Printf(format, v...)
	}
}

func (l *logger) Error(v ...interface{}) {
	if l.level <= ERROR {
		l.errorLogger.Println(v...)
	}
}

func (l *logger) Errorf(format string, v ...interface{}) {
	if l.level <= ERROR {
		l.errorLogger.Printf(format, v...)
	}
}

func (l *logger) Fatal(v ...interface{}) {
	l.fatalLogger.Println(v...)
	os.Exit(1)
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	l.fatalLogger.Printf(format, v...)
	os.Exit(1)
}
