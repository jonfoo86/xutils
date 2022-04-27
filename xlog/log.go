package xlog

import (
	"fmt"
	"io"
	"os"
	"runtime"
	goruntime "runtime"
	"sync"
	"time"

	"github.com/astaxie/beego/config"
)

var glogleve = 1
var glogfilepre = ""
var glogdir = "log"

var LogDebug = 1
var LogInfo = 2
var LogWarn = 3
var LogError = 4
var logoutput io.Writer = os.Stdout
var logoutputmutex sync.Mutex

// SetOutput sets the output destination for the logger.
func SetOutput(w io.Writer) {
	logoutputmutex.Lock()
	defer logoutputmutex.Unlock()
	logoutput = w
}

func GetOutput() io.Writer {
	logoutputmutex.Lock()
	defer logoutputmutex.Unlock()
	return logoutput
}

var once sync.Once
var bCreate bool = false

func init() {
	jsonconf, err := config.NewConfig("json", "./log.json")
	if err != nil {
		fmt.Println("cant load config error , msg :", err)
		fmt.Println("use default config")

	} else {
		glogleve, _ = jsonconf.Int("loglevel")
		glogfilepre = jsonconf.String("logfile")
	}

	////	 ErrorLevel =4    WarnLevel  =3   InfoLevel=2   DebugLevel=1

	lasthour := time.Now().Hour()
	if len(glogfilepre) == 0 {
		glogfilepre = string(os.Args[0]) + string(".log")
	}
	go func() {
		for {
			select {
			case <-time.After(time.Second):
				if lasthour != time.Now().Hour() {
					once = sync.Once{}
					lasthour = time.Now().Hour()
					bCreate = false
				}
			}
		}

	}()
}

var bCanCreate bool = false

func LogToFile() {
	bCanCreate = true
}

var startts int64 = time.Now().Unix()

func output(logptr *string) {

	if !bCreate {
		nowsec := time.Now().Unix()
		if bCanCreate || ((startts + int64(5)) < nowsec) {
			once.Do(func() { go CreateFile(glogdir, glogfilepre) })
			bCreate = true
		}
	}
	outpur := GetOutput()

	io.WriteString(outpur, *logptr)
}

func formatlog(levelInfo string, fileInfo bool, v ...interface{}) *string {
	header := ""
	now := time.Now().Format("2006-01-02 15:04:05")

	if fileInfo {
		var ok bool
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		}
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		header = fmt.Sprint(now, " ", file, ":", line, " ", levelInfo, " ")
	} else {
		header = fmt.Sprint(now, " ", levelInfo, " ")
	}

	lx := fmt.Sprint(v...)
	outstr := header + lx + "\n"
	return &outstr
}

func Debug(v ...interface{}) {
	if glogleve > LogDebug {
		return
	}
	s := formatlog("Debug", true, v...)

	switch goruntime.GOOS {
	case "darwin":
		output(s)

		break
	case "windows":
		fmt.Println(*s)
		break

	case "linux":
		output(s)
		break
	}

}

func Info(v ...interface{}) {
	if glogleve > LogInfo {
		return
	}
	s := formatlog("Info", true, v...)
	switch goruntime.GOOS {
	case "darwin":
		output(s)

		break
	case "windows":
		fmt.Println(*s)
		break

	case "linux":
		output(s)
		break
	}
}

func InfoWithoutFileInfo(v ...interface{}) {
	if glogleve > LogInfo {
		return
	}
	s := formatlog("Info", false, v...)
	switch goruntime.GOOS {
	case "darwin":
		output(s)

		break
	case "windows":
		fmt.Println(*s)
		break

	case "linux":
		output(s)
		break
	}
}

func Warn(v ...interface{}) {
	if glogleve > LogWarn {
		return
	}
	s := formatlog("Warn", true, v...)
	output(s)

	switch goruntime.GOOS {

	case "windows":
		fmt.Println(*s)
		break

	}
}

func Error(v ...interface{}) {
	if glogleve > LogError {
		return
	}
	s := formatlog("Error", true, v...)
	output(s)
	switch goruntime.GOOS {

	case "windows":
		fmt.Println(*s)
		break

	}
}

func Panic(v ...interface{}) {
	s := formatlog("Panic", true, v...)
	output(s)
	switch goruntime.GOOS {

	case "windows":
		fmt.Println(*s)
		break

	}
	time.Sleep(time.Millisecond * 300)
	panic(v)
}
