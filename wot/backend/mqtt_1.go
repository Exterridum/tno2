package backend

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/conas/tno2/util/sec"
	"github.com/conas/tno2/util/str"
	"github.com/conas/tno2/wot/server"
	"github.com/eclipse/paho.mqtt.golang"
)

//MQTT_1 is mqtt backend type 1
//type 1 mqtt backend supports single value properties, events and no actions
//type 1 mqtt backend is not conversation based
type MQTT_1 struct {
	client mqtt.Client
	values map[string]interface{}
}

func NewMQTT_1(cfg map[string]interface{}) Backend {
	url := cfg["url"].(string)
	username := cfg["username"].(string)
	password := cfg["password"].(string)

	id, _ := sec.UUID4()
	opts := mqtt.NewClientOptions().AddBroker(url).SetClientID(id)
	opts.SetKeepAlive(20 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetUsername(username)
	opts.SetPassword(password)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return &MQTT_1{
		client: c,
		values: make(map[string]interface{}),
	}
}

//TODO: Implement encoder
func (mb *MQTT_1) Bind(wos *server.WotServer, ctxPath string, encoder Encoder) {
	mb.setup(ctxPath, wos)
}

func (mb *MQTT_1) Start() {}

func (mb *MQTT_1) setup(ctxPath string, wos *server.WotServer) {
	deviceTopic := str.Concat(ctxPath, "/#")
	token2 := mb.client.Subscribe(deviceTopic, 0, mb.eventHandler(ctxPath, wos))
	if token2.Wait() && token2.Error() != nil {
		log.Fatal(token2.Error)
		os.Exit(1)
	}
	log.Info("MQTT_1 Backend: subscribed to device topic -> ", deviceTopic)

	for _, p := range wos.GetDescription().Properties {
		propPath := str.Concat(ctxPath, "/", p.Name)
		wos.OnGetProperty(p.Name, func() interface{} {
			return mb.values[propPath]
		})

		if p.Writable {
			wos.OnUpdateProperty(p.Name, func(payload interface{}) {
				mb.publish(propPath, payload)
			})
		}
	}
}

func (mb *MQTT_1) publish(propPath string, data interface{}) mqtt.Token {
	token := mb.client.Publish(propPath, 0, false, data)
	log.Info("Published ", propPath, " : ", data)
	return token
}

type PropertyChange struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

func (mb *MQTT_1) eventHandler(ctxPath string, wos *server.WotServer) func(mqtt.Client, mqtt.Message) {
	return func(client mqtt.Client, m mqtt.Message) {
		topic := m.Topic()

		p := PropertyChange{
			Name:  topic[len(ctxPath)+1 : len(topic)],
			Value: string(m.Payload()),
		}

		mb.values[topic] = p.Value
		wos.EmitEvent("property-change", p)
	}
}
