package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

var total = 0

// startSdrPower : start SDR Power Monitor
func startSdrPower(ctx context.Context) {
	timer := time.NewTicker(time.Second * time.Duration(syslogInterval))
	defer timer.Stop()
	log.Println("start sdr power")
	count := 0
	mHz := 24
	for {
		select {
		case <-timer.C:
			syslogCh <- fmt.Sprintf("type=Stats,total=%d,count=%d,ps=%.2f,send=%d,param=%d",
				total, count, float64(count)/float64(syslogInterval), syslogCount, sdr)
			log.Printf("type=Stats,total=%d,count=%d,ps=%.2f,send=%d,param=%d",
				total, count, float64(count)/float64(syslogInterval), syslogCount, sdr)
			syslogCount = 0
			count = 0
			mHz = 24
			sendMonitor()
		case <-ctx.Done():
			log.Println("stop sdr power")
			return
		default:
			if mHz < 1677 {
				doScan(mHz)
				count++
				mHz++
			}
		}

	}
}

// スキャンの実施
func doScan(mHz int) {

}

// syslogでレポートを送信する
func sendReport() {
}
