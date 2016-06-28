/*
 * Copyright (c) 2013 IBM Corp.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v1.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-v10.html
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 */

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/websocket"

	"github.com/eclipse/paho.mqtt.golang"
)

func main() {
	ch := mqttSubscribe("/foo/val", "tcp://localhost:1883", "gotrivial")
	http.Handle("/echo", websocket.Handler(newWsHandler(ch)))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func newWsHandler(ch chan []byte) func(*websocket.Conn) {
	return func(ws *websocket.Conn) {
		for {
			req := <-ch
			ws.Write(req)
		}
	}
}

func newMqttHandler(ch chan []byte) func(mqtt.Client, mqtt.Message) {
	return func(c mqtt.Client, m mqtt.Message) {
		ch <- m.Payload()
	}
}

func mqttSubscribe(topic, url, clientID string) chan []byte {
	opts := mqtt.NewClientOptions().AddBroker(url).SetClientID(clientID)
	opts.SetKeepAlive(20 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	ch := make(chan []byte)
	token := c.Subscribe(topic, 0, newMqttHandler(ch))

	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	return ch
}
