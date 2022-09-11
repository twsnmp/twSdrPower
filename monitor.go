package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	gopsnet "github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

var lastMonitorTime int64
var lastBytesRecv uint64
var lastBytesSent uint64

// sendMonitor : センサーが稼働するPCのリソース情報を送信する
func sendMonitor() {
	msg := "type=Monitor,"
	cpus, err := cpu.Percent(0, false)
	if err != nil {
		log.Printf("sendMonitor err=%v", err)
		return
	}
	msg += fmt.Sprintf("cpu=%.3f", cpus[0])
	loads, err := load.Avg()
	if err != nil {
		log.Printf("sendMonitor err=%v", err)
		return
	}
	msg += fmt.Sprintf(",load=%.3f", loads.Load1)
	mems, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("sendMonitor err=%v", err)
		return
	}
	msg += fmt.Sprintf(",mem=%.3f", mems.UsedPercent)
	nets, err := gopsnet.IOCounters(false)
	if err != nil {
		log.Printf("sendMonitor err=%v", err)
		return
	}
	now := time.Now().Unix()
	if lastMonitorTime > 0 {
		diff := now - lastMonitorTime
		if diff > 0 {
			dSent := nets[0].BytesSent - lastBytesSent
			dRecv := nets[0].BytesRecv - lastBytesRecv
			rxSpeed := 8.0 * float64(dRecv) / float64(diff)
			rxSpeed /= (1000 * 1000)
			txSpeed := 8.0 * float64(dSent) / float64(diff)
			txSpeed /= (1000 * 1000)
			msg += fmt.Sprintf(",recv=%d,sent=%d,rxSpeed=%.3f,txSpeed=%.3f",
				dRecv, dSent, rxSpeed, txSpeed)
		}
	}
	lastMonitorTime = time.Now().Unix()
	lastBytesRecv = nets[0].BytesRecv
	lastBytesSent = nets[0].BytesSent
	pids, err := process.Pids()
	if err != nil {
		log.Printf("sendMonitor err=%v", err)
		return
	}
	msg += fmt.Sprintf(",process=%d,param=%d", len(pids), sdr)
	syslogCh <- msg
}
