package backend

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/conas/tno2/util/async"
	"github.com/conas/tno2/util/col"
	"github.com/conas/tno2/util/sec"
	"github.com/conas/tno2/util/str"
	"github.com/conas/tno2/wot/server"
	"github.com/eclipse/paho.mqtt.golang"
)

type MQTT struct {
	client   mqtt.Client
	bindings map[string]*col.Map
}

func NewMQTT(cfg map[string]interface{}) Backend {
	url := cfg["url"].(string)
	opts := mqtt.NewClientOptions().AddBroker(url).SetClientID("ClientID")
	opts.SetKeepAlive(20 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return &MQTT{
		client:   c,
		bindings: make(map[string]*col.Map),
	}
}

func (mb *MQTT) Bind(wos *server.WotServer, baseTopic string, encoder Encoder) {
	bindingID, _ := sec.UUID4()
	mb.bindings[bindingID] = col.NewConcurentMap()

	mb.setupDeviceInTopic(bindingID, baseTopic, wos, encoder)
	mb.setupDeviceOutTopic(bindingID, baseTopic, wos, encoder)
}

func (mb *MQTT) Start() {}

func (mb *MQTT) setupDeviceInTopic(bindingID string, baseTopic string, wos *server.WotServer, encoder Encoder) {
	deviceInTopic := str.Concat(baseTopic, "/i")
	log.Info("MQTTBackend: device in topic -> ", deviceInTopic)

	for _, a := range wos.GetDescription().Actions {
		wos.OnInvokeAction(a.Name, func(payload interface{}, ph async.ProgressHandler) interface{} {
			log.Info("Action invoked ", a.Name, payload)
			return mb.publish(bindingID, encoder, deviceInTopic, BE_ACTION_RQ, a.Name, payload)
		})
	}

	for _, p := range wos.GetDescription().Properties {
		wos.OnGetProperty(p.Name, func() interface{} {
			return mb.publish(bindingID, encoder, deviceInTopic, BE_GET_PROP_RQ, p.Name, nil)
		})

		if p.Writable {
			wos.OnUpdateProperty(p.Name, func(payload interface{}) {
				mb.publish(bindingID, encoder, deviceInTopic, BE_SET_PROP_RQ, p.Name, payload)
			})
		}
	}
}

func (mb *MQTT) publish(
	bindingID string,
	encoder Encoder,
	deviceInTopic string,
	msgType int8,
	msgName string,
	data interface{}) interface{} {

	conversationID, _ := sec.UUID4()
	urlQ := encoder.Encode(msgType, conversationID, msgName, data)

	var response interface{}
	var promise *async.Promise
	if msgType == BE_ACTION_RQ || msgType == BE_GET_PROP_RQ {
		promise = async.NewPromise()
		mb.bindings[bindingID].Add(conversationID, promise)
	}

	log.Info("Will publish ", deviceInTopic, " : ", string(urlQ))
	mb.client.Publish(deviceInTopic, 0, false, urlQ)
	// wait to receive response on deviceOutTopic to fulfuill the promise
	// Q: should we timeout?
	if msgType == BE_ACTION_RQ || msgType == BE_GET_PROP_RQ {
		response = promise.Get()
		mb.bindings[bindingID].Del(conversationID)
	}

	return response
}

func (mb *MQTT) setupDeviceOutTopic(bindingID string, baseTopic string, wos *server.WotServer, encoder Encoder) {
	deviceOutTopic := str.Concat(baseTopic, "/o")
	log.Info("MQTTBackend: device out topic -> ", deviceOutTopic)
	token2 := mb.client.Subscribe(deviceOutTopic, 0, outSubHandler(wos, encoder, mb.bindings[bindingID]))
	if token2.Wait() && token2.Error() != nil {
		os.Exit(1)
	}
}

func outSubHandler(wos *server.WotServer, encoder Encoder, conversations *col.Map) func(mqtt.Client, mqtt.Message) {
	return func(client mqtt.Client, m mqtt.Message) {
		msgType, conversationID, msgName, msgData := encoder.Decode(m.Payload())

		log.Info("MQTT message receive ", string(m.Payload()))

		switch msgType {
		case BE_ACTION_RS:
			conv, _ := conversations.Get(conversationID)
			conv.(*async.Promise).Set(msgData)
		case BE_GET_PROP_RS:
			conv, _ := conversations.Get(conversationID)
			conv.(*async.Promise).Set(msgData)
		case BE_EVENT:
			wos.EmitEvent(msgName, msgData)
		}
	}
}
