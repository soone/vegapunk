package clog

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync/atomic"

	"github.com/soone/vegapunk/initialize"
	"github.com/soone/vegapunk/notification"
)

const (
	LDEBUG = iota + 1
	LWARN
	LINFO
	LERROR
	LFATAL
)

var Logx *CLog = Default()

type CLog struct {
	level        int64
	w            io.Writer
	debugLogger  *log.Logger
	warnLogger   *log.Logger
	infoLogger   *log.Logger
	errLogger    *log.Logger
	fatalLogger  *log.Logger
	notification map[string]notification.Notification
	prefix       string
}

func Default() *CLog {
	return New(os.Stderr, LDEBUG, "[clog]", 0)
}

func New(w io.Writer, level int64, prefix string, flag int) *CLog {
	if w == nil {
		w = os.Stderr
	}

	if level == 0 {
		level = LDEBUG
	}

	if prefix == "" {
		prefix = "CLOG"
	}

	if flag <= 0 {
		flag = log.LstdFlags | log.Lshortfile
	}

	return &CLog{
		level:       level,
		w:           w,
		debugLogger: log.New(w, fmt.Sprintf("[%s][DEBUG]", prefix), flag),
		warnLogger:  log.New(w, fmt.Sprintf("[%s][WARN]", prefix), flag),
		infoLogger:  log.New(w, fmt.Sprintf("[%s][INFO]", prefix), flag),
		errLogger:   log.New(w, fmt.Sprintf("[%s][ERROR]", prefix), flag),
		fatalLogger: log.New(w, fmt.Sprintf("[%s][FATAL]", prefix), flag),
		prefix:      prefix,
	}
}

func (l *CLog) RegistNotify(name string, notify notification.Notification) {
	if l.notification == nil {
		l.notification = make(map[string]notification.Notification)
	}

	l.notification[name] = notify
}

func (l *CLog) SetLevel(level int64) {
	if level < LDEBUG || level > LFATAL {
		level = LDEBUG
	}

	atomic.StoreInt64(&l.level, level)
}

func (l *CLog) Debugf(format string, v ...any) {
	if atomic.LoadInt64(&l.level) > LDEBUG {
		return
	}

	l.debugLogger.Printf(format, v...)
}

func (l *CLog) Debugln(v ...any) {
	if atomic.LoadInt64(&l.level) > LDEBUG {
		return
	}

	l.debugLogger.Println(v...)
}

func (l *CLog) NDebugf(format string, v ...any) {
	l.Debugf(format, v...)
	l.sendNotify(fmt.Sprintf(format, v...))
}

func (l *CLog) NDebugln(v ...any) {
	l.Debugln(v...)
	l.sendNotify(fmt.Sprintln(v...))
}

func (l *CLog) Warnf(format string, v ...any) {
	if atomic.LoadInt64(&l.level) > LWARN {
		return
	}

	l.warnLogger.Printf(format, v...)
}

func (l *CLog) Warnln(v ...any) {
	if atomic.LoadInt64(&l.level) > LWARN {
		return
	}

	l.warnLogger.Println(v...)
}

func (l *CLog) NWarnf(format string, v ...any) {
	l.Warnf(format, v...)
	l.sendNotify(fmt.Sprintf(format, v...))
}

func (l *CLog) NWarnln(v ...any) {
	l.Warnln(v...)
	l.sendNotify(fmt.Sprintln(v...))
}

func (l *CLog) Infof(format string, v ...any) {
	if atomic.LoadInt64(&l.level) > LINFO {
		return
	}

	l.infoLogger.Printf(format, v...)
}

func (l *CLog) Infoln(v ...any) {
	if atomic.LoadInt64(&l.level) > LINFO {
		return
	}

	l.infoLogger.Println(v...)
}

func (l *CLog) NInfof(format string, v ...any) {
	l.Infof(format, v...)
	l.sendNotify(fmt.Sprintf(format, v...))
}

func (l *CLog) NInfoln(v ...any) {
	l.Infoln(v...)
	l.sendNotify(fmt.Sprintln(v...))
}

func (l *CLog) Errorf(format string, v ...any) {
	if atomic.LoadInt64(&l.level) > LERROR {
		return
	}

	l.errLogger.Printf(format, v...)
}

func (l *CLog) Errorln(v ...any) {
	if atomic.LoadInt64(&l.level) > LERROR {
		return
	}

	l.errLogger.Println(v...)
}

func (l *CLog) NErrorf(format string, v ...any) {
	l.Errorf(format, v...)
	l.sendNotify(fmt.Sprintf(format, v...))
}

func (l *CLog) NErrorln(v ...any) {
	l.Errorln(v...)
	l.sendNotify(fmt.Sprintln(v...))
}

func (l *CLog) Fatalf(format string, v ...any) {
	if atomic.LoadInt64(&l.level) > LFATAL {
		return
	}

	l.fatalLogger.Printf(format, v...)
}

func (l *CLog) Fatalln(v ...any) {
	if atomic.LoadInt64(&l.level) > LFATAL {
		return
	}

	l.fatalLogger.Println(v...)
}

func (l *CLog) NFatalf(format string, v ...any) {
	l.Fatalf(format, v...)
	l.sendNotify(fmt.Sprintf(format, v...))
}

func (l *CLog) NFatalln(v ...any) {
	l.Fatalln(v...)
	l.sendNotify(fmt.Sprintln(v...))
}

func (l *CLog) sendNotify(msg string) {
	if l.notification == nil {
		return
	}

	msg = fmt.Sprintf("[%s]%s", l.prefix, msg)

	for name, notify := range l.notification {
		initialize.WG2Exec(func(args ...any) {
			name := args[0].(string)
			notify := args[1].(notification.Notification)

			err := notify.Send(msg)
			if err != nil {
				l.Errorf("[SEND-NOTIFY][%s] error: %v", name, err)
			}
		}, name, notify)
	}
}
