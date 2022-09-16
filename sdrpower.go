package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/samuel/go-rtlsdr/rtl"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var total = 0
var scan = 0
var dev *rtl.Device
var id = time.Now().Unix()
var data = []opts.LineData{}

// startSdrPower : start SDR Power Monitor
func startSdrPower(ctx context.Context) {
	log.Println("start sdr power")

	timer := time.NewTicker(time.Second * time.Duration(syslogInterval))
	defer timer.Stop()
	count := 0
	hz := startHz
	dur := int64(0)
	for {
		select {
		case <-timer.C:
			if hz >= endHz {
				hz = startHz
				id = time.Now().Unix()
				scan++
			}
			syslogCh <- fmt.Sprintf("type=Stats,total=%d,count=%d,ps=%.2f,send=%d,param=%d,scan=%d,dur=%d",
				total, count, float64(count)/float64(syslogInterval), syslogCount, sdr, scan, dur)
			log.Printf("type=Stats,total=%d,count=%d,ps=%.2f,send=%d,param=%d,scan=%d,dur=%d",
				total, count, float64(count)/float64(syslogInterval), syslogCount, sdr, scan, dur)
			syslogCount = 0
			count = 0
			sendMonitor()
		case <-ctx.Done():
			if dev != nil {
				dev.Close()
				dev = nil
			}
			log.Println("stop sdr power")
			return
		default:
			if hz <= endHz {
				if dev == nil {
					if err := openRTLSdr(); err != nil {
						log.Printf("failed to open RTL-SDR err=%v", err)
						hz = endHz + stepHz
						dur = -1
						continue
					}
					log.Println("open RTL-SDR")
				}
				if err := doScan(hz); err != nil {
					log.Printf("failed to scan err=%v", err)
					hz = endHz + stepHz
					dur = -1
					continue
				}
				count++
				total++
				hz += stepHz
				if hz > endHz {
					dur = time.Now().Unix() - id
					dev.Close()
					dev = nil
					log.Println("close RTL-SDR")
					outChart()
				}
			} else {
				time.Sleep(time.Millisecond * 100)
			}
		}

	}
}

// スキャンの実施
func doScan(hz int) error {
	// set center freq
	if err := dev.SetCenterFreq(uint(hz)); err != nil {
		return err
	}
	// wait 5mSec
	time.Sleep(time.Millisecond * 5)
	// Dummy Read
	dmy := make([]byte, 1<<12)
	if _, err := dev.Read(dmy); err != nil {
		return err
	}
	// read data
	buf := make([]byte, 16384)
	n, err := dev.Read(buf)
	if err != nil {
		return err
	}
	p := 0
	t := 0
	for i := 0; i < n; i++ {
		s := int(buf[i]) - 127
		t += s
		p += (s * s)
	}
	dc := float64(t) / float64(n)
	e := float64(t)*2.0*dc - dc*dc*float64(n)
	p -= int(math.Round(e))
	dbm := float64(p)
	dbm /= float64(1e6) // 1M
	dbm = 10 * math.Log10(dbm)
	data = append(data, opts.LineData{Value: []float64{float64(hz) / 1e6, dbm}})
	syslogCh <- fmt.Sprintf("type=Power,id=%x,freq=%d,dbm=%.3f", id, hz, dbm)
	return nil
}

func openRTLSdr() error {
	var err error
	dev, err = rtl.Open(sdr)
	if err != nil {
		return err
	}
	// no direct sample
	// no offset tuning
	// set auto gain
	if gain == 0 {
		if err := dev.SetTunerGainMode(false); err != nil {
			return err
		}
	} else {
		if err := dev.SetTunerGainMode(true); err != nil {
			return err
		}
		if err := dev.SetTunerGain(gain); err != nil {
			return err
		}
	}
	// no PPM
	// disable biasTee
	if err := dev.SetBiasTee(false); err != nil {
		return err
	}
	// reset buffer
	if err := dev.ResetBuffer(); err != nil {
		return err
	}
	// set sample rate 1MHz
	if err := dev.SetSampleRate(1_000_000); err != nil {
		return err
	}
	if err := dev.SetCenterFreq(uint(24 * 1e6)); err != nil {
		return err
	}
	// wait 100mSec
	time.Sleep(time.Millisecond * 100)
	// read dumy data
	buf := make([]byte, 16384)
	if _, err := dev.Read(buf); err != nil {
		return err
	}
	return nil
}

func outChart() {
	if chartTitle == "" || len(data) < 1 {
		data = []opts.LineData{}
		return
	}
	theme := "white"
	if dark {
		theme = "dark"
	}
	title := chartTitle + "-" + time.Now().Format("2006/01/02/ 15:04:05")
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: title}),
		charts.WithInitializationOpts(opts.Initialization{Theme: theme}),
		charts.WithDataZoomOpts(opts.DataZoom{Type: "slider"}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true, Trigger: "axis"}),
		charts.WithXAxisOpts(opts.XAxis{
			AxisLabel: &opts.AxisLabel{Show: true, Formatter: "{value} MHz"},
			Type:      "value",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			AxisLabel: &opts.AxisLabel{Show: true, Formatter: "{value} dbm"},
			Type:      "value",
		}),
	)
	line.AddSeries("power", data).SetSeriesOptions(
		charts.WithLineChartOpts(opts.LineChart{
			Smooth: true,
		}),
	)
	file := filepath.Join(chartFolder, chartTitle+"-"+time.Now().Format("20060102150405")+".html")
	f, err := os.Create(file)
	if err != nil {
		log.Printf("chart save err=%v", err)
	} else {
		line.Render(f)
		defer f.Close()
	}
	data = []opts.LineData{}
}
