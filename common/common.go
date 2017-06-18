package common

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/kisom/goutils/fileutil"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/lestrrat/go-strftime"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

const DefaultDirPerm = 0740

func InitLogrus(name, level, dir string, pid bool, rotate bool) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Fatalf("not supported log level: %s", level)
	} else {
		logrus.SetLevel(lvl)
	}

	if err := os.MkdirAll(dir, DefaultDirPerm); err != nil {
		return
	}

	link := path.Join(dir, fmt.Sprintf("%s.log", name))
	fileName := name
	if pid {
		fileName = fmt.Sprintf("%s-%d", name, os.Getpid())
	}
	p := path.Join(dir, fileName) + ".%Y%m%d%H%M"
	logrus.Debugf("log pattern: %q", p)

	var writer io.Writer
	if rotate {
		options := []rotatelogs.Option{
			rotatelogs.WithLinkName(link),
			rotatelogs.WithLocation(time.UTC),
		}
		if rotate {
			options = append(options, rotatelogs.WithRotationTime(time.Hour))
		}
		writer, err = rotatelogs.New(
			p,
			options...,
		)
	} else {
		strfobj, err := strftime.New(p)
		if err != nil {
			logrus.Fatalf("problem with parsing object: %s", err)
			//return nil, errors.Wrap(err, `invalid strftime pattern`)
		}
		p = strfobj.FormatString(time.Now().UTC())
		logrus.Infof("log: %q", p)
		writer, err = os.Create(p)

		if err == nil {
			logrus.Debugf("setup link: %q", link)
			if fileutil.FileDoesExist(link) {
				if err = os.Remove(link); err != nil {
					logrus.Debugf("link exist, remove")
				}
			}
			err = os.Symlink(p, link)
		}
	}
	if err != nil {
		logrus.Fatalf("cannot setup logging: %s", err)
	}

	m := lfshook.WriterMap{}

	for _, v := range logrus.AllLevels {
		if v <= lvl {
			m[v] = writer
		}
	}
	logrus.AddHook(lfshook.NewHook(m))
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(os.Stderr)
	logrus.Debug("file hook added")
}
