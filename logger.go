package main

import (
	"fmt"
	"os"
	"strings"
	"time"
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

// Trace 跟踪日志
func Trace(format string, v ...interface{}) {
	if currentLogLevel <= TRACE {
		printLog("TRACE", format, v...)
	}
}

// Debug 调试日志
func Debug(format string, v ...interface{}) {
	if currentLogLevel <= DEBUG {
		printLog("DEBUG", format, v...)
	}
}

// Info 信息日志
func Info(format string, v ...interface{}) {
	if currentLogLevel <= INFO {
		printLog("INFO", format, v...)
	}
}

// Warn 警告日志
func Warn(format string, v ...interface{}) {
	if currentLogLevel <= WARN {
		printLog("WARN", format, v...)
	}
}

// Error 错误日志
func Error(format string, v ...interface{}) {
	if currentLogLevel <= ERROR {
		printLog("ERROR", format, v...)
	}
}

// Fatal 致命错误日志
func Fatal(format string, v ...interface{}) {
	printLog("FATAL", format, v...)
	fmt.Println("\n按回车键退出...")
	fmt.Scanln()
	os.Exit(1)
}

// printLog 内部日志打印函数
func printLog(level string, format string, v ...interface{}) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf(format, v...)
	fmt.Printf("[%s] %s %s\n", level, timestamp, message)
}
