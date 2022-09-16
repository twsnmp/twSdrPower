package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/samuel/go-rtlsdr/rtl"
)

var version = "v1.0.0"
var commit = ""
var syslogDst = ""

var sdr = 0
var startHzStr = "24M"
var endHzStr = "1667M"
var stepHzStr = "1M"
var startHz = 0
var endHz = 0
var stepHz = 0
var gain = 0
var chartFolder = "./"
var chartTitle = ""
var dark = false
var list = false
var once = false
var syslogInterval = 600

func init() {
	flag.StringVar(&syslogDst, "syslog", "", "syslog destnation list")
	flag.IntVar(&sdr, "sdr", 0, "RTL-SDR Device Number")
	flag.IntVar(&gain, "gain", 0, "RTL-SDR Tuner gain (0=auto)")
	flag.StringVar(&startHzStr, "start", "24M", "start frequency")
	flag.StringVar(&endHzStr, "end", "1667M", "end frequency")
	flag.StringVar(&stepHzStr, "step", "1M", "step frequency")
	flag.StringVar(&chartTitle, "chart", "", "chart title")
	flag.StringVar(&chartFolder, "folder", "./", "chart folder")
	flag.BoolVar(&dark, "dark", false, "dark mode chart")
	flag.BoolVar(&list, "list", false, "List RTL-STR")
	flag.BoolVar(&once, "once", false, "Only once")
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
	if list {
		showDevices()
		return
	}
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
	setScanRange()
	log.Printf("version=%s", fmt.Sprintf("%s(%s)", version, commit))
	log.Printf("sdr=%d,chart=%s,interval=%d,start=%d,end=%d,step=%d,gain=%d", sdr, chartTitle, syslogInterval, startHz, endHz, stepHz, gain)
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

func setScanRange() {
	var err error
	startHz, err = getHz(startHzStr)
	if err != nil {
		log.Fatalf("setRange start err=%v", err)
	}
	endHz, err = getHz(endHzStr)
	if err != nil {
		log.Fatalf("setRange end err=%v", err)
	}
	stepHz, err = getHz(stepHzStr)
	if err != nil {
		log.Fatalf("setRange step err=%v", err)
	}
	if startHz < 24*1e6 {
		log.Fatalf("setRange  start < 24MHz")
	}
	if endHz > 1667*1e6 {
		log.Fatalf("setRange  end < 1667MHz")
	}
	if startHz >= endHz {
		log.Fatalf("setRange  start > end")
	}
	if stepHz == 0 {
		log.Fatalf("setRange  step == 0")
	}
	bin := (endHz - startHz) / stepHz
	if bin > 2000 {
		log.Fatalf("setRange  bin > 2000")
	}
}

// Get frequency
func getHz(s string) (int, error) {
	if len(s) < 1 {
		return 0, fmt.Errorf("no frequency")
	}
	last := s[len(s)-1:]
	suff := float64(1.0)
	switch last {
	case "g", "G":
		suff *= 1e3
		fallthrough
	case "m", "M":
		suff *= 1e3
		fallthrough
	case "k", "K":
		suff *= 1e3
		f, err := strconv.ParseFloat(s[0:len(s)-1], 64)
		if err != nil {
			return 0, err
		}
		return int(f * suff), nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return int(f), nil
}

func showDevices() {
	count := rtl.DeviceCount()
	fmt.Printf("Device List count=%d\n", count)
	for i := 0; i < count; i++ {
		name := rtl.DeviceName(i)
		m, p, sn, err := rtl.DeviceUSBStrings(i)
		if err != nil {
			log.Printf("DeviceUSBStrings Failed: %v", err)
			continue
		}
		fmt.Printf("%d,%s,%s,%s,%s\n", i, name, m, p, sn)
	}
}
