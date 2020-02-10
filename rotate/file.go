package rotate

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bingoohuang/gometrics/util"
	"github.com/sirupsen/logrus"
)

const yyyyMMdd = "yyyy-MM-dd"

type File struct {
	Filename   string
	MaxBackups int

	lastDay string
	dir     string
	file    *os.File
}

// NewFile create a rotation option
func NewFile(filename string, maxBackups int) (*File, error) {
	o := &File{
		Filename:   filename,
		MaxBackups: maxBackups,
		dir:        filepath.Dir(filename),
	}

	if err := os.MkdirAll(o.dir, 0755); err != nil {
		return nil, err
	}

	if err := o.open(); err != nil {
		return nil, err
	}

	return o, nil
}

func (o *File) open() error {
	f, err := os.OpenFile(o.Filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logrus.Warnf("log file %s created error %v", o.Filename, err)

		return err
	}

	o.file = f

	logrus.Infof("log file %s created", o.Filename)

	return nil
}

func (o *File) rotateFiles(t time.Time) error {
	rotated, outMaxBackups := o.detectRotate(t)

	return o.doRotate(rotated, outMaxBackups)
}

func (o *File) doRotate(rotated string, outMaxBackups []string) error {
	if rotated != "" {
		if err := o.close(); err != nil {
			return err
		}

		if err := os.Rename(o.Filename, rotated); err != nil {
			logrus.Infof("rotate %s to %s error %v", o.Filename, rotated, err)
			return err
		}

		logrus.Infof("%s rotated to %s", o.Filename, rotated)

		if err := o.open(); err != nil {
			return err
		}
	}

	for _, old := range outMaxBackups {
		if err := os.Remove(old); err != nil {
			logrus.Warnf("remove log file %s before max backup days %d error %v", old, o.MaxBackups, err)
		}

		logrus.Infof("%s before max backup days %d removed", old, o.MaxBackups)
	}

	return nil
}

func (o *File) close() error {
	if o.file == nil {
		return nil
	}

	err := o.file.Close()
	o.file = nil

	return err
}

func (o *File) detectRotate(t time.Time) (rotated string, outMaxBackups []string) {
	ts := util.FormatTime(t, yyyyMMdd)

	if o.lastDay == "" {
		o.lastDay = ts
	}

	if o.lastDay != ts {
		o.lastDay = ts

		yesterday := t.AddDate(0, 0, -1)
		rotated = o.Filename + "." + util.FormatTime(yesterday, yyyyMMdd)
	}

	if o.MaxBackups > 0 {
		day := t.AddDate(0, 0, -o.MaxBackups)
		_ = filepath.Walk(o.dir, func(path string, fi os.FileInfo, err error) error {
			if err != nil || fi.IsDir() {
				return err
			}

			if strings.HasPrefix(path, o.Filename+".") {
				fis := path[len(o.Filename+"."):]
				if backDay, err := util.ParseTime(fis, yyyyMMdd); err != nil {
					return nil // ignore this file
				} else if backDay.Before(day) {
					outMaxBackups = append(outMaxBackups, path)
				}
			}

			return nil
		})
	}

	return rotated, outMaxBackups
}

// Close closes the file
func (o *File) Close() error {
	return o.close()
}

func (o *File) write(d []byte, flush bool) (int, error) {
	if err := o.rotateFiles(time.Now()); err != nil {
		return 0, err
	}

	n, err := o.file.Write(d)

	if err != nil {
		return n, err
	}

	if flush {
		err = o.file.Sync()
	}

	return n, err
}

// Write writes data to a file
func (o *File) Write(d []byte) (int, error) {
	return o.write(d, false)
}

// Flush flushes the file
func (o *File) Flush() error {
	return o.file.Sync()
}
