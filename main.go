package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var config struct {
	Verbose   bool
	Interval  int
	ProbePort int
	CU        string
	CT        string
	CM        string
}

func main() {
	var flags struct {
		Server   string
		Port     int
		User     string
		Password string
	}

	flag.StringVar(&flags.Server, "server", "", "server address")
	flag.StringVar(&flags.User, "user", "", "user")
	flag.IntVar(&flags.Port, "port", 35601, "port")
	flag.StringVar(&flags.Password, "password", "USER_DEFAULT_PASSWORD", "user")
	flag.IntVar(&config.Interval, "interval", 1, "update interval(s)")
	flag.IntVar(&config.ProbePort, "probe-port", 80, "probe port")
	flag.StringVar(&config.CU, "probe-cu-host", "cu.tz.cloudcpp.com", "China Unicom probe host")
	flag.StringVar(&config.CT, "probe-ct-host", "ct.tz.cloudcpp.com", "China Telecom probe host")
	flag.StringVar(&config.CM, "probe-cm-host", "cm.tz.cloudcpp.com", "China Mobile probe host")
	flag.BoolVar(&config.Verbose, "verbose", false, "verbose mode")
	flag.Parse()

	if flags.Server == "" || flags.User == "" {
		flag.Usage()
		return
	}

	addr := fmt.Sprintf("%s:%d", flags.Server, flags.Port)
	sigCh := make(chan os.Signal, 1)

	go setupDaemon(addr, flags.Password, flags.User)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigCh
}
