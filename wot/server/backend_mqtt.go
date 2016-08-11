package server

import (
	"os"
	"strings"
	"time"

	"github.com/conas/tno2/util/str"
	"github.com/eclipse/paho.mqtt.golang"
)

type MQTT struct {
	client mqtt.Client
}

func MQTTBackend(url string) *MQTT {
	opts := mqtt.NewClientOptions().AddBroker(url).SetClientID("ClientID")
	opts.SetKeepAlive(20 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return &MQTT{
		client: c,
	}
}

func (mb *MQTT) Bind(baseTopic string, s *WotServer) {
	s.OnGetProperty

	// inTopic := str.Concat(baseTopic, "/i")
	// token1 := mb.client.Subscribe(inTopic, 0, inHandler(s))
	// if token1.Wait() && token1.Error() != nil {
	// 	os.Exit(1)
	// }

	outTopic := str.Concat(baseTopic, "/o")
	token2 := mb.client.Subscribe(outTopic, 0, outHandler(s))
	if token2.Wait() && token2.Error() != nil {
		os.Exit(1)
	}
}

//QUICK POC IMPL

func outHandler(s *WotServer) func(mqtt.Client, mqtt.Message) {
	return func(c mqtt.Client, m mqtt.Message) {
		p := string(m.Payload())
		nd := strings.Split(p, ":")
		eventName := nd[0]
		eventData := parseEvent(nd[1])
		s.EmitEvent(eventName, eventData)
	}
}

func parseEvent(payload string) map[string]string {
	kvs := strings.Split(payload, "&")

	m := make(map[string]string, len(kvs))

	for _, s := range kvs {
		kv := strings.Split(s, "=")
		m[kv[0]] = kv[1]
	}

	return m
}
