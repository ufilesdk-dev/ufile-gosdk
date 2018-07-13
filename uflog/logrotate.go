/*=============================================================================
#     FileName: logrotate.go
#         Desc: log rotate
#       Author: ato.ye
#        Email: ato.ye@ucloud.cn
#     HomePage: http://www.ucloud.cn
#      Version: 0.0.1
#   LastChange: 2016-02-3 20:23:20
#      History:
=============================================================================*/
package uflog

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var isOnOpenFlag int32 = 0

type RotateLogger struct {
	dir         string
	prefix      string
	suffix      string
	size        int64
	cnt         uint64
	checkRotate chan uint64
	isClosed    chan bool
	idx         int64
	logf        *os.File
	level       int
	*Logger
}

func NewRotateLogger(dir, prefix, suffix string, size int64) (rl *RotateLogger, err error) {

	if dir == "" {
		dir = "."
	}
	if dir[len(dir)-1:] == "/" {
		dir = dir[:len(dir)-1]
	}

	var filelists []string
	if filelists, err = filepath.Glob(TodayLogPrefix(dir, prefix) + "*"); err != nil {
		return
	}

	logfile := ""
	var idx, mtime int64 = 0, 0
	// 判断文件的最近修改时间，选择最新的文件作为日志文件
	for _, fname := range filelists {
		if f, err := os.Open(fname); err == nil {
			fi, _ := f.Stat()
			if fi.ModTime().Unix() >= mtime {
				mtime = fi.ModTime().Unix()
				logfile = fname
			}
			f.Close()
		}
	}
	// 判断文件的大小,和设置的大小对比，如果文件>设置的大小,索引+1
	if logfile != "" {
		if i, err := ParseIdx(logfile, dir, prefix, suffix); err == nil {
			idx = i
		}
		f, _ := os.Open(logfile)
		fi, _ := f.Stat()
		if fi.Size() > size {
			idx += 1
		}
		f.Close()
	}

	rl = &RotateLogger{
		dir:         dir,
		prefix:      prefix,
		suffix:      suffix,
		size:        size,
		cnt:         0,
		checkRotate: make(chan uint64),
		isClosed:    make(chan bool),
		idx:         idx,
	}
	var logf *os.File
	if logf, err = rl.OpenLogf(); err != nil {
		return
	}
	rl.logf = logf
	rl.SetLogf(rl.logf)
	go lorateLog(rl)
	return
}

func (rl *RotateLogger) SetLogf(f *os.File) {
	rl.logf = f
	if rl.Logger != nil {
		rl.Logger.isClosed <- true
	}
	rl.Logger = New(rl.logf, "", Ldefault)
	rl.Logger.SetOutputLevel(rl.level)
}

func (rl *RotateLogger) SetOutputLevel(lvl int) {
	rl.level = lvl
	rl.Logger.SetOutputLevel(lvl)
}

func (rl *RotateLogger) Close() {
	rl.isClosed <- true
}

func (rl *RotateLogger) OpenLogf() (logf *os.File, err error) {
	if err = os.MkdirAll(rl.dir, 0777); err != nil {
		return
	}
	fname := LogName(rl.dir, rl.prefix, rl.suffix, rl.idx)
	logf, err = os.OpenFile(fname, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return
	}
	return
}

func (rl *RotateLogger) Output(reqId string, lvl int, calldepth int, s string) error {
	rl.checkRotate <- atomic.AddUint64(&rl.cnt, 1)
	return rl.Logger.Output(reqId, lvl, calldepth, s)
}

func lorateLog(rl *RotateLogger) {
	for {
		select {
		case cnt := <-rl.checkRotate:
			if cnt%100 == 0 {
				if fi, err := rl.logf.Stat(); err == nil {
					if fi.Size() >= rl.size {
						rl.idx += 1
						logf, err := rl.OpenLogf()
						if err == nil {
							rl.SetLogf(logf)
						}
					}
				}
			}
		case isClosed := <-rl.isClosed:
			if isClosed {
				rl.Logger.isClosed <- true
				return
			}
		}
	}
}

//====================================================================================================

func LogName(dir, prefix, suffix string, idx int64) string {

	t := time.Now()
	return fmt.Sprintf("%s/%s%d%02d%02d%d%s", dir, prefix, t.Year(), t.Month(), t.Day(), idx, suffix)
}

func TodayLogPrefix(dir, prefix string) string {

	t := time.Now()
	return fmt.Sprintf("%s/%s%d%02d%02d", dir, prefix, t.Year(), t.Month(), t.Day())
}

func ParseIdx(filename, dir, prefix, suffix string) (idx int64, err error) {

	if dir != "" {
		pos := -1
		if pos = strings.Index(filename, dir); pos != -1 {
			filename = filename[len(dir)+1:]
		} else {
			return -1, errors.New("invalid log filename")
		}
	}

	re, _ := regexp.Compile(prefix + "[0-9]{8}(.*)" + suffix)
	matchList := re.FindStringSubmatch(filename)
	if len(matchList) != 2 {
		return -1, errors.New("parse idx failed: invalid format")
	}
	idx, err = strconv.ParseInt(matchList[1], 10, 32)
	return
}
