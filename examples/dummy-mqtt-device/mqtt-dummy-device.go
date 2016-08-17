package dmqttdevice

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/conas/tno2/util/str"
	"github.com/eclipse/paho.mqtt.golang"
)

func Setup(url, deviceTopic string) {
	opts := mqtt.NewClientOptions().AddBroker(url).SetClientID("ClientID")
	opts.SetKeepAlive(20 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	deviceOutTopic := str.Concat(deviceTopic, "/o")
	deviceInTopic := str.Concat(deviceTopic, "/i")
	token2 := c.Subscribe(deviceInTopic, 0, func(client mqtt.Client, m mqtt.Message) {
		log.Info("MQTT Dummy - Message received on: ", deviceOutTopic, " data: ", string(m.Payload()))
	})

	if token2.Wait() && token2.Error() != nil {
		os.Exit(1)
	}
}
