package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttCh = make(chan interface{}, 2000)

type mqttPowerDataEnt struct {
	Time string    `json:"time"`
	ID   string    `json:"id"`
	Freq []int     `json:"freq"`
	DBM  []float64 `json:"dbm"`
}

type mqttSdrStatsDataEnt struct {
	Time  string  `json:"time"`
	Total int     `json:"total"`
	Count int     `json:"count"`
	PS    float64 `json:"ps"`
	Send  int     `json:"send"`
	SDR   int     `json:"sdr"`
	Scan  int     `json:"scan"`
	Dur   int64   `json:"dur"`
}

type mqttMonitorDataEnt struct {
	Time    string  `json:"time"`
	CPU     float64 `json:"cpu"`
	Memory  float64 `json:"memory"`
	Load    float64 `json:"load"`
	Sent    uint64  `json:"sent"`
	Recv    uint64  `json:"recv"`
	TxSpeed float64 `json:"tx_speed"`
	RxSpeed float64 `json:"rx_speed"`
	Process int     `json:"process"`
}

func startMQTT(ctx context.Context) {
	if mqttDst == "" {
		return
	}
	broker := mqttDst
	if !strings.Contains(broker, "://") {
		broker = "tcp://" + broker
	}
	if strings.LastIndex(broker, ":") <= 5 {
		broker += ":1883"
	}
	log.Printf("start mqtt broker=%s", broker)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	if mqttUser != "" && mqttPassword != "" {
		opts.SetUsername(mqttUser)
		opts.SetPassword(mqttPassword)
	}
	opts.SetClientID(mqttClientID)
	opts.SetAutoReconnect(true)
	if debug {
		opts.OnConnect = connectHandler
		opts.OnConnectionLost = connectLostHandler
	}
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		return
	}
	defer client.Disconnect(250)
	for {
		select {
		case <-ctx.Done():
			log.Println("stop mqtt")
			return
		case msg := <-mqttCh:
			if s := makeMqttData(msg); s != "" {
				if debug {
					log.Println(s)
				}
				client.Publish(getMqttTopic(msg), 1, false, s).Wait()
			}
		}
	}
}

func getMqttTopic(msg interface{}) string {
	r := mqttTopic
	switch msg.(type) {
	case *mqttMonitorDataEnt:
		r += "/Monitor"
	case *mqttPowerDataEnt:
		r += "/Power"
	case *mqttSdrStatsDataEnt:
		r += "/Stats"
	default:
		log.Printf("getMqttTopic: unknown msg type %T", msg)
	}
	return r
}

func makeMqttData(msg interface{}) string {
	if j, err := json.Marshal(msg); err == nil {
		return string(j)
	}
	return ""
}

func publishMQTT(msg interface{}) {
	if mqttDst == "" {
		return
	}
	select {
	case mqttCh <- msg:
	default:
		if debug {
			log.Println("mqtt channel full, skipping message")
		}
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
}
