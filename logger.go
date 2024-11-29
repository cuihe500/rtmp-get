package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// LogLevel 日志级别
type LogLevel int

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
)

var (
	// 当前日志级别
	currentLogLevel = INFO
	// 调试模式
	debugMode = false
)

// Logger 自定义日志结构
type Logger struct {
	level LogLevel
	*log.Logger
}

var (
	traceLogger = &Logger{TRACE, log.New(os.Stdout, "[TRACE] ", log.LstdFlags)}
	debugLogger = &Logger{DEBUG, log.New(os.Stdout, "[DEBUG] ", log.LstdFlags)}
	infoLogger  = &Logger{INFO, log.New(os.Stdout, "[INFO] ", log.LstdFlags)}
	warnLogger  = &Logger{WARN, log.New(os.Stdout, "[WARN] ", log.LstdFlags)}
	errorLogger = &Logger{ERROR, log.New(os.Stdout, "[ERROR] ", log.LstdFlags)}
)

// SetLogLevel 设置日志级别
func SetLogLevel(level string) {
	switch strings.ToUpper(level) {
	case "TRACE":
		currentLogLevel = TRACE
	case "DEBUG":
		currentLogLevel = DEBUG
	case "INFO":
		currentLogLevel = INFO
	case "WARN":
		currentLogLevel = WARN
	case "ERROR":
		currentLogLevel = ERROR
	default:
		currentLogLevel = INFO
	}
	Debug("日志级别设置为: %s", level)
}

// SetDebugMode 设置调试模式
func SetDebugMode(debug bool) {
	debugMode = debug
	if debug {
		currentLogLevel = DEBUG // 在调试模式下强制设置为 DEBUG 级别
		Debug("调试模式已启用")
	}
}

// IsDebugMode 获取当前是否为调试模式
func IsDebugMode() bool {
	return debugMode
}

func (l *Logger) log(format string, v ...interface{}) {
	if l.level >= currentLogLevel {
		if len(v) == 0 {
			l.Println(format)
		} else {
			l.Printf(format, v...)
		}
	}
}

// Trace 跟踪日志
func Trace(format string, v ...interface{}) {
	if debugMode || currentLogLevel <= TRACE {
		traceLogger.log(format, v...)
	}
}

// Debug 调试日志
func Debug(format string, v ...interface{}) {
	if debugMode || currentLogLevel <= DEBUG {
		debugLogger.log(format, v...)
	}
}

// Info 信息日志
func Info(format string, v ...interface{}) {
	if currentLogLevel <= INFO {
		infoLogger.log(format, v...)
	}
}

// Warn 警告日志
func Warn(format string, v ...interface{}) {
	if currentLogLevel <= WARN {
		warnLogger.log(format, v...)
	}
}

// Error 错误日志
func Error(format string, v ...interface{}) {
	if currentLogLevel <= ERROR {
		errorLogger.log(format, v...)
	}
}

// Fatal 致命错误日志
func Fatal(format string, v ...interface{}) {
	errorLogger.log(format, v...)
	fmt.Println("\n按回车键退出...")
	fmt.Scanln()
	os.Exit(1)
}
