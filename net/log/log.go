// Package log has the same logging functions as google.golang.org/appengine/log
// but allows to change
package log

import (
	"github.com/golang/glog"
	"golang.org/x/net/context"
)

type LogOverrideFunc func(level int64, format string, args ...interface{})

var logOverrideKey = "holds a LogOverrideFunc"

func WithLogOverride(ctx context.Context, f LogOverrideFunc) context.Context {
	return context.WithValue(ctx, &logOverrideKey, f)
}

func logf(ctx context.Context, level int64, format string, args ...interface{}) {
	// if f, ok := ctx.Value(&logOverrideKey).(logOverrideFunc); ok {
	// 	f(level, format, args...)
	// 	return
	// }
	// logf(fromContext(ctx), level, format, args...)
	panic("not implemented")
}

func gen(lvl int64) func(ctx context.Context, format string, args ...interface{}) {
	return func(ctx context.Context, format string, args ...interface{}) {
		logf(ctx, lvl, format, args...)
	}
}

func Wrap(f func(format string, args ...interface{})) func(ctx context.Context, format string, args ...interface{}) {
	return func(ctx context.Context, format string, args ...interface{}) {
		f(format, args...)
	}
}

var (
	Debugf    LogFunc = Wrap(glog.Infof)
	Infof     LogFunc = Wrap(glog.Infof)
	Warningf  LogFunc = Wrap(glog.Warningf)
	Errorf    LogFunc = Wrap(glog.Errorf)
	Criticalf LogFunc = Wrap(glog.Fatalf)
)

type LogFunc func(ctx context.Context, format string, args ...interface{})
