package log

import (
	"io"
	"log"
	"os"
)

var (
	l *dsLogger = nil
)

const (
	LevelFatal = iota
	LevelError = iota
	LevelInfo  = iota
	LevelDebug = iota
)

const (
	prefixDebug = "[DEBUG]"
	prefixInfo  = "[INFO]"
	prefixError = "[ERROR]"
	prefixFatal = "[FATAL]"
)

type dsLogger struct {
	level  int
	debugL *log.Logger
	infoL  *log.Logger
	errorL *log.Logger
	fatalL *log.Logger
}

func InitLog(prefix string, logfileName string, level int) {
	if l != nil {
		return
	}
	l = new(dsLogger)
	l.level = level
	mw := io.MultiWriter(os.Stdout)
	if logfileName != "" {
		f, err := os.Open(logfileName)
		if err != nil {
			log.Fatal(err)
		}
		mw = io.MultiWriter(os.Stdout, f)
	}
	if level >= LevelDebug {
		l.debugL = log.New(mw, prefix+prefixDebug, log.LstdFlags)
	}
	if level >= LevelInfo {
		l.infoL = log.New(mw, prefix+prefixInfo, log.LstdFlags)
	}
	if level >= LevelError {
		l.errorL = log.New(mw, prefix+prefixError, log.LstdFlags)
	}
	if level >= LevelFatal {
		l.fatalL = log.New(mw, prefix+prefixFatal, log.LstdFlags)
	}
}

func GetInfoLogger() *log.Logger {
	return l.infoL
}
func Debug(v ...interface{}) {
	if l != nil && l.debugL != nil {
		l.debugL.Println(v...)
	} else {
		log.Println(v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if l != nil && l.debugL != nil {
		l.debugL.Printf(format, v...)
	} else {
		log.Printf(format, v...)
	}
}

func Info(v ...interface{}) {
	if l != nil && l.infoL != nil {
		l.infoL.Println(v...)
	} else {
		log.Println(v...)
	}
}

func Infof(format string, v ...interface{}) {
	if l != nil && l.infoL != nil {
		l.infoL.Printf(format, v...)
	} else {
		log.Printf(format, v...)
	}
}

func Error(v ...interface{}) {
	if l != nil && l.errorL != nil {
		l.errorL.Println(v...)
	} else {
		log.Println(v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if l != nil && l.errorL != nil {
		l.errorL.Printf(format, v...)
	} else {
		log.Printf(format, v...)
	}
}

func Fatal(v ...interface{}) {
	if l != nil && l.fatalL != nil {
		l.fatalL.Println(v...)
	} else {
		log.Println(v...)
	}
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	if l != nil && l.fatalL != nil {
		l.fatalL.Printf(format, v...)
	} else {
		log.Printf(format, v...)
	}
	os.Exit(1)
}
