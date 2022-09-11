package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var syslogCh chan string
var syslogCount = 0

func startSyslog(ctx context.Context) {
	syslogCh = make(chan string, 2000)
	dstList := strings.Split(syslogDst, ",")
	dst := []net.Conn{}
	for _, d := range dstList {
		if !strings.Contains(d, ":") {
			d += ":514"
		}
		s, err := net.Dial("udp", d)
		if err != nil {
			log.Fatal(err)
		}
		syslogCh <- fmt.Sprintf("start send syslog to %s", d)
		dst = append(dst, s)
	}
	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}
	defer func() {
		for _, d := range dst {
			d.Close()
		}
	}()
	for {
		select {
		case <-ctx.Done():
			log.Println("stop syslog")
			return
		case msg := <-syslogCh:
			syslogCount++
			s := fmt.Sprintf("<%d>%s %s twWifiScan: %s", 21*8+6, time.Now().Format("2006-01-02T15:04:05-07:00"), host, msg)
			for _, d := range dst {
				d.Write([]byte(s))
			}
		}
	}
}
