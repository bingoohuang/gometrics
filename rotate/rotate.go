package rotate

import (
	"os"
	"path/filepath"
	"time"

	"github.com/bingoohuang/gometrics/util"
)

// File describes a file that gets rotating
type File struct {
	// info about currently opened file
	Path       string
	MaxBackups int

	lastTime string
	file     *os.File
}

func (f *File) close() error {
	if f.file == nil {
		return nil
	}

	err := f.file.Close()
	f.file = nil

	return err
}

func (f *File) open() error {
	// we can't assume that the dir for the file already exists
	dir := filepath.Dir(f.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var err error

	if f.file, err = os.OpenFile(f.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err != nil {
		return err
	}

	return nil
}

const yyyyMMdd = "yyyy-MM-dd"

// rotate if required
func (f *File) rotate() error {
	if f.file == nil {
		_ = f.open()
	}

	t := time.Now()
	ts := util.FormatTime(t, yyyyMMdd)

	if f.lastTime == "" {
		f.lastTime = ts
		return nil
	}

	if f.lastTime == ts {
		return nil
	}

	f.lastTime = ts

	if err := f.close(); err != nil {
		return err
	}

	yesterday := t.AddDate(0, 0, -1)
	if err := os.Rename(f.Path, f.Path+"."+util.FormatTime(yesterday, yyyyMMdd)); err != nil {
		return err
	}

	if f.MaxBackups > 0 {
		day := t.AddDate(0, 0, -f.MaxBackups)
		_ = os.Remove(f.Path + "." + util.FormatTime(day, yyyyMMdd))
	}

	return f.open()
}

// NewFile creates a new file that will be rotated daily (at midnight in specified location).
// logPath is file full path like /var/log/my.log
func NewFile(logPath string, maxBackups int) (*File, error) {
	f := &File{Path: logPath, MaxBackups: maxBackups}

	// force early failure if we can't open the file
	if err := f.rotate(); err != nil {
		return nil, err
	}

	return f, nil
}

// Close closes the file
func (f *File) Close() error {
	return f.close()
}

func (f *File) write(d []byte, flush bool) (int, error) {
	if err := f.rotate(); err != nil {
		return 0, err
	}

	n, err := f.file.Write(d)

	if err != nil {
		return n, err
	}

	if flush {
		err = f.file.Sync()
	}

	return n, err
}

// Write writes data to a file
func (f *File) Write(d []byte) (int, error) {
	return f.write(d, false)
}

// Flush flushes the file
func (f *File) Flush() error {
	return f.file.Sync()
}
