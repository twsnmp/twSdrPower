package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var version = "v1.0.0"
var commit = ""
var syslogDst = ""
var sdr = 0
var chartFolder = "./"
var chartTitle = ""
var dark = false
var syslogInterval = 600

func init() {
	flag.StringVar(&syslogDst, "syslog", "", "syslog destnation list")
	flag.StringVar(&chartTitle, "chart", "", "chart title")
	flag.StringVar(&chartFolder, "folder", "./", "chart folder")
	flag.BoolVar(&dark, "dark", false, "dark mode chart")
	flag.IntVar(&sdr, "sdr", 0, "RTL-SDR Device Number")
	flag.IntVar(&syslogInterval, "interval", 600, "syslog send interval(sec)")
	flag.VisitAll(func(f *flag.Flag) {
		if s := os.Getenv("TWSDRPOWER_" + strings.ToUpper(f.Name)); s != "" {
			f.Value.Set(s)
		}
	})
	flag.Parse()
}

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().Format("2006-01-02T15:04:05.999 ") + string(bytes))
}

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
	log.Printf("version=%s", fmt.Sprintf("%s(%s)", version, commit))
	log.Printf("sdr=%d,chart=%s,interval=%d", sdr, chartTitle, syslogInterval)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go startSyslog(ctx)
	go startSdrPower(ctx)
	<-quit
	syslogCh <- "quit by signal"
	time.Sleep(time.Second * 1)
	log.Println("quit by signal")
	cancel()
	time.Sleep(time.Second * 2)
}
