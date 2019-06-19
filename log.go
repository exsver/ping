package ping

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	LogLevel *LogSeverity
)

type LogSeverity struct {
	Fail    *log.Logger
	Info    *log.Logger
	Message *log.Logger
}

func (logSeverity *LogSeverity) Default() {
	LogLevel = &LogSeverity{
		Info:    log.New(ioutil.Discard, "[Verbose]: ", log.LstdFlags|log.Lshortfile),
		Fail:    log.New(os.Stderr, "[Error]: ", log.LstdFlags|log.Lshortfile),
		Message: log.New(ioutil.Discard, "[Debug]: ", log.LstdFlags|log.Lshortfile),
	}
}

func (logSeverity *LogSeverity) Verbose() {
	LogLevel = &LogSeverity{
		Info:    log.New(os.Stdout, "[Verbose]: ", log.LstdFlags|log.Lshortfile),
		Fail:    log.New(os.Stderr, "[Error]: ", log.LstdFlags|log.Lshortfile),
		Message: log.New(ioutil.Discard, "[Debug]: ", log.LstdFlags|log.Lshortfile),
	}
}

func (logSeverity *LogSeverity) Debug() {
	LogLevel = &LogSeverity{
		Info:    log.New(os.Stdout, "[Verbose]: ", log.LstdFlags|log.Lshortfile),
		Fail:    log.New(os.Stderr, "[Error]: ", log.LstdFlags|log.Lshortfile),
		Message: log.New(os.Stdout, "[Debug]: ", log.LstdFlags|log.Lshortfile),
	}
}

func (logSeverity *LogSeverity) Silent() {
	LogLevel = &LogSeverity{
		Info:    log.New(ioutil.Discard, "[Verbose]: ", log.LstdFlags|log.Lshortfile),
		Fail:    log.New(ioutil.Discard, "[Error]: ", log.LstdFlags|log.Lshortfile),
		Message: log.New(ioutil.Discard, "[Debug]: ", log.LstdFlags|log.Lshortfile),
	}
}
