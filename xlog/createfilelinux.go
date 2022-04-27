//go:build !windows
// +build !windows

package xlog

import (
	"os"
	"strconv"
	"syscall"
	"time"
)

func CreateFile(dir, filename string) {
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)
	os.MkdirAll(dir, os.FileMode(0755))
	os.Remove(filename)

	realfile := filename + "." + time.Now().Format("2006-01-02_15:04:05") + "." + strconv.Itoa(os.Getpid())
	file, err := os.OpenFile(dir+"/"+realfile, os.O_WRONLY|os.O_CREATE, os.FileMode(0666))
	if err == nil {
		SetOutput(file)
	} else {
		SetOutput(os.Stdout)
	}

	os.Symlink(dir+"/"+realfile, filename)
}
