package main

import (
	"github.com/getsentry/sentry-go"
	"log"
	"time"

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
		Sentry   bool
	}

	flag.StringVar(&flags.Server, "server", "", "server address")
	flag.StringVar(&flags.User, "user", "", "user")
	flag.IntVar(&flags.Port, "port", 35601, "port")
	flag.StringVar(&flags.Password, "password", "USER_DEFAULT_PASSWORD", "user")
	flag.BoolVar(&flags.Sentry, "sentry", false, "use Sentry")
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

	if flags.Sentry {
		err := sentry.Init(sentry.ClientOptions{
			Dsn: "https://7a776e797a044b34b6db9bc4905b3437@o45713.ingest.sentry.io/5310005",
		})
		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}
		defer sentry.Flush(2 * time.Second)
	}

	addr := fmt.Sprintf("%s:%d", flags.Server, flags.Port)
	sigCh := make(chan os.Signal, 1)

	go setupDaemon(addr, flags.Password, flags.User)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigCh
}
