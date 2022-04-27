package xdeamon

import (
	"flag"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/juju/fslock"

	"github.com/jonfoo86/xutils/xlog"

	"github.com/sevlyar/go-daemon"
)

var (
	signal = flag.String("s", "", `send signal to the daemon
		quit — graceful shutdown
		stop — fast shutdown
		reload — reloading the configuration file`)
)

var endfunc func() = nil

var lock *fslock.Lock = nil

func Deam(runfunc func(), clearfunc func(), svrname string) {
	flag.Parse()
	daemon.AddCommand(daemon.StringFlag(signal, "quit"), syscall.SIGQUIT, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "reload"), syscall.SIGHUP, reloadHandler)

	cntxt := &daemon.Context{
		PidFileName: "pid." + svrname,
		PidFilePerm: 0644,
		LogFileName: "stdout." + svrname + ".log",
		LogFilePerm: 0644,
		WorkDir:     "./",
		Umask:       027,
		Args:        os.Args,
	}

	endfunc = clearfunc

	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			fmt.Println("Unable send signal to the daemon:", err)
			return
		}
		daemon.SendCommands(d)
		return
	} else {
		d, err := cntxt.Search()
		if err != nil {
		} else {
			err := d.Signal(syscall.Signal(0))
			if err == nil {
				fmt.Println("old deam is alive ")
				return
			} else {
				//fmt.Printf("Search Deam Resp: ", err.Error())
			}
		}
	}

	d, err := cntxt.Reborn()
	if err != nil {
		fmt.Println("Unable Reborn:", err)
		return
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	lock = fslock.New("lock." + svrname)
	lockErr := lock.TryLock()
	if lockErr != nil {
		fmt.Println("falied to acquire lock > " + lockErr.Error())
		return
	}

	fmt.Println("got the lock")
	//lock file
	xlog.LogToFile()
	xlog.Info("svr run with daemon --->>>>>>>>>")

	go runfunc()

	time.Sleep(time.Second * 5)

	err = daemon.ServeSignals()
	if err != nil {
		xlog.Info("Error:", err)
	}
	xlog.Info("daemon terminated")
}

func termHandler(sig os.Signal) error {
	xlog.Info("terminating...")
	go func() {
		select {
		case <-time.After(time.Millisecond * 3500):
			xlog.Info("force exit...")
			time.Sleep(time.Millisecond * 50)
			os.Exit(0)
		}
	}()
	if endfunc != nil {
		endfunc()
	}
	time.Sleep(time.Millisecond * 500)
	return daemon.ErrStop
}

func reloadHandler(sig os.Signal) error {
	xlog.Info("configuration reloaded")
	return nil
}
