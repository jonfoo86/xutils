package main

import (
	"github.com/jonfoo86/xutils/xlog"
)

func main() {
	xlog.Info("test", " test 2 ")
	xlog.LogToFile()
	xlog.Info("test", " test 2 ")

}
