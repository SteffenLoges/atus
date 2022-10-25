package logger

import (
	"atus/backend/sqlite"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

var OnLog = func(LogEntry, LogSeverity, ...string) {}

type LogEntry struct {
	Time         time.Time `json:"time"`
	UID          string    `json:"uid"`
	RefType      RefType   `json:"refType"`
	forceConsole bool
	LogType      LogType     `json:"logType"`
	LogSeverity  LogSeverity `json:"logSeverity"`
	Message      string      `json:"message"`
}

func new() *LogEntry {
	return &LogEntry{
		Time:    time.Now(),
		LogType: TypeGeneric,
	}
}

type RefType string

const (
	RefRelease    RefType = "RELEASE"
	RefFileserver RefType = "FILESERVER"
	RefSource     RefType = "SOURCE"
)

func Ref(refType RefType, uid string) *LogEntry {
	return &LogEntry{
		Time:    time.Now(),
		UID:     uid,
		RefType: refType,
		LogType: TypeGeneric,
	}
}

func ForceConsole() *LogEntry {
	le := new()
	le.forceConsole = true
	return le
}

func (le *LogEntry) ForceConsole() *LogEntry {
	le.forceConsole = true
	return le
}

type LogType string

const (
	TypeGeneric    LogType = "GENERIC"
	TypeFileserver LogType = "FILESERVER"
	TypePredb      LogType = "PREDB"
	TypeSample     LogType = "SAMPLE"
	TypeRelease    LogType = "RELEASE"
	TypeSource     LogType = "SOURCE"
)

func Type(logType LogType) *LogEntry {
	le := new()
	le.LogType = logType
	return le
}

func (le *LogEntry) Type(logType LogType) *LogEntry {
	le.LogType = logType
	return le
}

type LogSeverity int

const (
	SeverityDebug LogSeverity = iota
	SeverityInfo
	SeverityWarning
	SeverityError
	SeverityFatal
)

// -- debug -------------------------------------

func Debug(message ...string) {
	new().log(SeverityDebug, message...)
}

func Debugf(format string, a ...interface{}) {
	new().log(SeverityDebug, fmt.Sprintf(format, a...))
}

func (le *LogEntry) Debug(message ...string) {
	le.log(SeverityDebug, message...)
}

func (le *LogEntry) Debugf(format string, a ...interface{}) {
	le.log(SeverityDebug, fmt.Sprintf(format, a...))
}

// -- info --------------------------------------

func Info(message ...string) {
	new().log(SeverityInfo, message...)
}

func Infof(format string, a ...interface{}) {
	new().log(SeverityInfo, fmt.Sprintf(format, a...))
}

func (le *LogEntry) Info(message ...string) {
	le.log(SeverityInfo, message...)
}

func (le *LogEntry) Infof(format string, a ...interface{}) {
	le.log(SeverityInfo, fmt.Sprintf(format, a...))
}

// -- warning -----------------------------------

func Warning(message ...string) {
	new().log(SeverityWarning, message...)
}

func Warningf(format string, a ...interface{}) {
	new().log(SeverityWarning, fmt.Sprintf(format, a...))
}

func (le *LogEntry) Warning(message ...string) {
	le.log(SeverityWarning, message...)
}

func (le *LogEntry) Warningf(format string, a ...interface{}) {
	le.log(SeverityWarning, fmt.Sprintf(format, a...))
}

// -- error -------------------------------------

func Error(message ...string) {
	new().log(SeverityError, message...)
}

func Errorf(format string, a ...interface{}) {
	new().log(SeverityError, fmt.Sprintf(format, a...))
}

func (le *LogEntry) Error(message ...string) {
	le.log(SeverityError, message...)
}

func (le *LogEntry) Errorf(format string, a ...interface{}) {
	le.log(SeverityError, fmt.Sprintf(format, a...))
}

// -- fatal -------------------------------------

func Fatal(message ...string) {
	new().log(SeverityFatal, message...)
}

func Fatalf(format string, a ...interface{}) {
	new().log(SeverityFatal, fmt.Sprintf(format, a...))
}

func (le *LogEntry) Fatal(message ...string) {
	le.log(SeverityFatal, message...)
}

func (le *LogEntry) Fatalf(format string, a ...interface{}) {
	le.log(SeverityFatal, fmt.Sprintf(format, a...))
}

// ----------------------------------------------

var LogCache []*LogEntry
var m sync.RWMutex

func addToCache(le *LogEntry) {
	m.Lock()
	defer m.Unlock()

	// keep a max of 250 entries
	if len(LogCache) >= 250 {
		LogCache = LogCache[1:]
	}

	LogCache = append(LogCache, le)
}

func GetCache() []*LogEntry {
	m.RLock()
	defer m.RUnlock()

	return LogCache
}

func (le *LogEntry) log(severity LogSeverity, message ...string) {

	if le.UID == "" {
		le.UID = "N/A"
	}

	le.LogSeverity = severity
	le.Message = strings.Join(message, " ")

	if severity == SeverityError || le.forceConsole {
		if severity == SeverityError {
			color.New(color.FgRed, color.Bold).Println(le.Time.Format(time.UnixDate), "|", severity, "|", le.LogType, "|", le.UID, "|", le.Message)
		} else {
			fmt.Println(le.Time.Format(time.UnixDate), "|", severity, "|", le.LogType, "|", le.UID, "|", le.Message)
		}
	}

	if severity != SeverityDebug {
		go le.writeLog(severity, message...)
	}

	OnLog(*le, severity, message...)

	if severity == SeverityFatal {
		log.Fatal(le.Message)
	}

	addToCache(le)

}

func (le *LogEntry) writeLog(severity LogSeverity, message ...string) {

	_, err := sqlite.Conn.Exec(
		`INSERT INTO
				log
				(
					severity,
					type,
					message,
					release_uid
				)	VALUES (?,?,?,?)`,
		severity,
		le.LogType,
		strings.Join(message, " "),
		le.UID,
	)

	if err != nil {
		panic(err)
	}

}
