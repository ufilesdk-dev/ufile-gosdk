/*=============================================================================
#     FileName: api.go
#         Desc: log api
#       Author: ato.ye
#        Email: ato.ye@ucloud.cn
#     HomePage: http://www.ucloud.cn
#      Version: 0.0.1
#   LastChange: 2016-02-3 20:22:42
#      History:
=============================================================================*/
package uflog

import (
	"fmt"
	"os"
)

const (
	defaultLogDir         = "ufilelogs"
	defaultLogPrefix      = "ufilelog_"
	defaultLogSuffix      = ".log"
	defaultLogSize        = 50 // MB
	defaultLogLevelString = "DEBUG"
)

var Glogger *Logger

func InitLogger(dir, prefix, suffix string, size int64, level string) {
	if Glogger != nil {
		Glogger.Close()
	}
	if dir == "" {
		dir = defaultLogDir
	}
	if prefix == "" {
		prefix = defaultLogPrefix
	}
	if suffix == "" {
		suffix = defaultLogSuffix
	}
	if size <= 0 {
		size = defaultLogSize
	}
	if level == "" {
		level = defaultLogLevelString
	}
	logger, err := NewRotate(dir, prefix, suffix, size)
	if err != nil {
		fmt.Println("Init Logger fail:", err)
		os.Exit(-1)
	}
	Glogger = logger
	SetLogLevel(level)
}

func InitDefaultLogger() {
	InitLogger(defaultLogDir, defaultLogPrefix, defaultLogSuffix, defaultLogSize, defaultLogLevelString)
}

func SetLogLevel(level string) {
	if Glogger == nil {
		InitLogger(defaultLogDir, defaultLogPrefix, defaultLogSuffix, defaultLogSize, level)
	}
	Glogger.SetOutputLevelString(level)
}

func INFOF(format string, v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Infof(format, v...)
}

func INFO(v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Info(v...)
}

func ERRORF(format string, v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Errorf(format, v...)
}

func ERROR(v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Error(v...)
}

func WARN(v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Warn(v...)
}

func WARNF(format string, v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Warnf(format, v...)
}

func DEBUG(v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Debug(v...)
}

func DEBUGF(format string, v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Debugf(format, v...)
}
