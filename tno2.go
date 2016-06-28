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
	"log"
	"os"
	"time"

	"golang.org/x/net/websocket"

	"github.com/eclipse/paho.mqtt.golang"
)

func echoHandler(ws *websocket.Conn) {
	msg := make([]byte, 512)
	n, err := ws.Read(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Receive: %s\n", msg[:n])

	m, err := ws.Write(msg[:n])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Send: %s\n", msg[:m])
}

func newMqttHandler(ch chan string) func(mqtt.Client, mqtt.Message) {
	return func(c mqtt.Client, m mqtt.Message) {
		ch <- string(m.Payload())
	}
}

func mqttSubscribe(topic, url, clientID string) chan string {
	opts := mqtt.NewClientOptions().AddBroker(url).SetClientID(clientID)
	opts.SetKeepAlive(20 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	ch := make(chan string)
	token := c.Subscribe(topic, 0, newMqttHandler(ch))

	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	return ch
}

func main() {
	//websocket
	// http.Handle("/echo", websocket.Handler(echoHandler))
	// fmt.Printf("Listening 1 ...")
	// err := http.ListenAndServe(":8080", nil)
	// fmt.Printf("Listening 2 ...")
	// if err != nil {
	// 	panic("ListenAndServe: " + err.Error())
	// }

	ch := mqttSubscribe("/foo/val", "tcp://localhost:1883", "gotrivial")

	for {
		req := <-ch
		fmt.Printf("MQTT received: %s\n", req)
	}

	// c.Disconnect(250)

	// time.Sleep(1 * time.Second)
}
