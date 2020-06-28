package main

import (
	tcping "github.com/cloverstd/tcping/ping"

	"sync"
	"time"
)

var lrLock sync.Mutex
var lostRate = map[string]float64{
	"10010": 0.0,
	"189":   0.0,
	"10086": 0.0,
}
var ptLock sync.Mutex
var pingTime = map[string]int{
	"10010": 0,
	"189":   0,
	"10086": 0,
}

func startProbe() {
	go startPing(config.CM, "10086")
	time.Sleep(time.Second)
	go startPing(config.CT, "189")
	time.Sleep(time.Second)
	go startPing(config.CU, "10010")
}

func startPing(host, mark string) {
	lostPacket := 0
	allPacket := 0
	startTime := time.Now()

	for {
		pingTarget := &tcping.Target{
			Host:     host,
			Port:     config.ProbePort,
			Protocol: tcping.HTTP,
			Counter:  1,
			Interval: time.Second,
			Timeout:  time.Second * 10,
		}
		httpPing := tcping.NewTCPing()
		httpPing.SetTarget(pingTarget)

		ch := httpPing.Start()
		<-ch

		result := httpPing.Result()
		allPacket += result.Counter
		lostPacket += result.Counter - result.SuccessCounter
		ptLock.Lock()
		pingTime[mark] = int(result.TotalDuration)
		ptLock.Unlock()

		if allPacket > 60 {
			lrLock.Lock()
			lostRate[mark] = float64(lostPacket) / float64(allPacket)
			lrLock.Unlock()
		}

		if time.Since(startTime) > 3600*1000 {
			lostPacket = 0
			allPacket = 0
			startTime = time.Now()
		}

		time.Sleep(time.Second * 10)
	}
}
