//go:build windows
// +build windows

package xlog

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func CreateFile(dir, filename string) {

	os.MkdirAll(dir, os.FileMode(0755))
	os.Remove(filename)

	realfile := filename + "." + time.Now().Format("2006-01-02_15_04_05") + "." + strconv.Itoa(os.Getpid())
	file, err := os.Create(dir + "/" + realfile) //os.OpenFile(dir+"/"+realfile, os.O_WRONLY|os.O_CREATE, os.FileMode(0666))
	if err == nil {
		SetOutput(file)
	} else {

		fmt.Println("Fail Create LogFile :", err.Error())

		SetOutput(os.Stdout)
	}

	os.Symlink(dir+"/"+realfile, filename)
}
