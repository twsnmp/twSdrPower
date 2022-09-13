package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/samuel/go-rtlsdr/rtl"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var total = 0

// startSdrPower : start SDR Power Monitor
func startSdrPower(ctx context.Context) {
	dev, err := rtl.Open(0)
	if err != nil {
		log.Fatalf("rtl.Open Failed to open device err=%v", err)
	}
	defer dev.Close()
	// no direct sample
	// no offset tuning
	// set auto gain
	if err := dev.SetTunerGainMode(false); err != nil {
		log.Fatalf("dev.SetTunerGainMode err=%v", err)
	}
	// no PPM
	// disable biasTee
	if err := dev.SetBiasTee(false); err != nil {
		log.Fatalf("dev.SetBiasTee err=%v", err)
	}
	// reset buffer
	if err := dev.ResetBuffer(); err != nil {
		log.Fatalf("dev.ResetBuffer err=%v", err)
	}
	// set sample rate 1MHz
	if err := dev.SetSampleRate(1_000_000); err != nil {
		log.Fatalf("dev.ResetBuffer err=%v", err)
	}

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

func outChart(xaxis []int, data []opts.LineData) {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Freq MAP by RTL-SDR"}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show:  true,
			Right: "10%",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  true,
					Type:  "png",
					Title: "Save to Image",
				},
				DataZoom: &opts.ToolBoxFeatureDataZoom{
					Show: true,
				},
			}},
		),
	)
	line.SetXAxis(xaxis).
		AddSeries("power", data)
	f, _ := os.Create("power1.html")
	line.Render(f)
	defer f.Close()
}
