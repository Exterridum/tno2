package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/conas/tno2/util/async"
	"github.com/conas/tno2/util/sec"
	"github.com/conas/tno2/util/str"
	"github.com/conas/tno2/wot/encoder"
	"github.com/conas/tno2/wot/model"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// FIXME:
// 1. Too much repetititve code regarding codec handling

type Http struct {
	port          int
	router        *mux.Router
	hrefs         []string
	wotServers    map[string]*WotServer
	subscribers   *Subscribers
	actionResults *ActionResults
}

// ----- Server API methods

var hostname = "http://localhost:8080"

func NewHttp(port int) *Http {
	// r.PathPrefix("/model").HandlerFunc(Model)
	return &Http{
		port:          port,
		router:        mux.NewRouter().StrictSlash(true),
		hrefs:         make([]string, 0),
		wotServers:    make(map[string]*WotServer),
		subscribers:   NewSubscribers(),
		actionResults: NewActionResults(),
	}
}

func (p *Http) Bind(ctxPath string, s *WotServer) {
	td := s.GetDescription()
	p.wotServers[ctxPath] = s
	p.createRoutes(ctxPath, td)
	td.Uris = append(td.Uris, str.Concat(hostname, ctxPath))
}

func (p *Http) Start() {
	p.registerRoot()
	port := str.Concat(":", strconv.Itoa(p.port))
	log.Fatal(http.ListenAndServe(port, p.router))
}

// ----- ThingDescription parser methods

func (p *Http) createRoutes(ctxPath string, td *model.ThingDescription) {
	p.registerDeviceRoot(ctxPath, td)
	p.registerDeviceDescriptor(ctxPath, td)
	p.registerProperties(ctxPath, td)
	p.registerActions(ctxPath, td)
	p.registerEvents(ctxPath, td)
}

func (p *Http) registerRoot() {
	p.addRoute(&route{
		name:    "index",
		method:  "GET",
		pattern: "/",
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
			ls := links()

			for path := range p.wotServers {
				ls.Links = append(ls.Links, httpSubUrl(r, path))
			}

			sendJsonOK(w, ls)
		},
	})
}

func (p *Http) registerDeviceRoot(ctxPath string, td *model.ThingDescription) {
	p.addRoute(&route{
		name:    td.Name,
		method:  "GET",
		pattern: contextPath(ctxPath, ""),
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
			codec, err := p.getCodec(ctxPath, encoder.ENCODING_JSON)

			if err != nil {
				sendJsonERR(w, err)
				return
			}

			hrefs := links(httpSubUrl(r, "description"))

			sendOK(w, codec, hrefs)
		},
	})
}

func (p *Http) registerDeviceDescriptor(ctxPath string, td *model.ThingDescription) {
	p.addRoute(&route{
		name:    str.Concat(td.Name, "-descriptor"),
		method:  "GET",
		pattern: contextPath(ctxPath, "description"),
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
			codec, err := p.getCodec(ctxPath, encoder.ENCODING_JSON)

			if err != nil {
				sendJsonERR(w, err)
				return
			}

			sendOK(w, codec, td)
		},
	})
}

func (p *Http) registerProperties(ctxPath string, td *model.ThingDescription) {
	for _, prop := range td.Properties {
		p.addRoute(&route{
			name:        prop.Name,
			method:      "GET",
			pattern:     contextPath(ctxPath, prop.Hrefs[0]),
			handlerFunc: p.propertyGetHandler(ctxPath, &prop),
		})

		if prop.Writable {
			p.addRoute(&route{
				name:        prop.Hrefs[0],
				method:      "PUT",
				pattern:     contextPath(ctxPath, prop.Hrefs[0]),
				handlerFunc: p.propertySetHandler(ctxPath, &prop),
			})
		}

		prop.Hrefs[0] = str.Concat(hostname, ctxPath, "/", prop.Hrefs[0])
	}
}

func (p *Http) registerActions(ctxPath string, td *model.ThingDescription) {
	for _, action := range td.Actions {
		p.addRoute(&route{
			name:        action.Hrefs[0],
			method:      "POST",
			pattern:     contextPath(ctxPath, action.Hrefs[0]),
			handlerFunc: p.actionStartHandler(p.wotServers[ctxPath], action.Name),
		})

		p.addRoute(&route{
			name:        str.Concat(action.Hrefs[0], "Task-Rest"),
			method:      "GET",
			pattern:     contextPath(ctxPath, str.Concat(action.Hrefs[0], "/{taskid}")),
			handlerFunc: p.actionTaskHandler(p.wotServers[ctxPath]),
		})

		p.addRoute(&route{
			name:        str.Concat(action.Hrefs[0], "Task-WS"),
			method:      "GET",
			pattern:     contextPath(ctxPath, str.Concat(action.Hrefs[0], "/ws/{taskid}")),
			handlerFunc: p.actionWSTaskHandler(p.wotServers[ctxPath]),
		})

		action.Hrefs[0] = str.Concat(hostname, ctxPath, "/", action.Hrefs[0])
	}
}

func (p *Http) registerEvents(ctxPath string, td *model.ThingDescription) {
	for _, event := range td.Events {
		p.addRoute(&route{
			name:        event.Hrefs[0],
			method:      "POST",
			pattern:     contextPath(ctxPath, event.Hrefs[0]),
			handlerFunc: p.eventSubscribeHandler(p.wotServers[ctxPath], event.Name),
		})

		p.addRoute(&route{
			name:        str.Concat(event.Hrefs[0], "WebSocket"),
			method:      "GET",
			pattern:     contextPath(ctxPath, str.Concat(event.Hrefs[0], "/ws/{subscriptionID}")),
			handlerFunc: p.eventWSClientHandler(p.wotServers[ctxPath]),
		})

		event.Hrefs[0] = str.Concat(hostname, ctxPath, "/", event.Hrefs[0])
	}
}

type WotObject struct {
	Value interface{} `json:"value"`
}

func (w *WotObject) GetValue() interface{} {
	return w.GetValue()
}

func (p *Http) propertyGetHandler(ctxPath string, prop *model.Property) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		encoder, err := p.getCodec(ctxPath, encoder.ENCODING_JSON)

		if err != nil {
			sendJsonERR(w, err)
			return
		}

		promise, rc := p.wotServers[ctxPath].GetProperty(prop.Name)

		if rc == WOT_OK {
			value := promise.Get()
			sendOK(w, encoder, value)
		} else {
			sendERR(w, encoder, rc)
		}
	}
}

func (p *Http) propertySetHandler(ctxPath string, prop *model.Property) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		encoder, err := p.getCodec(ctxPath, encoder.ENCODING_JSON)

		if err != nil {
			sendJsonERR(w, err)
			return
		}

		wo := &WotObject{}
		err = encoder.Unmarshal(r.Body, wo)

		if err != nil {
			sendERR(w, encoder, err)
			return
		}

		promise, rc := p.wotServers[ctxPath].SetProperty(prop.Name, wo.GetValue())

		if rc == WOT_OK {
			promise.Get()
		} else {
			sendERR(w, encoder, rc)
		}
	}
}

func (p *Http) getCodec(ctxPath string, encoding encoder.Encoding) (encoder.Encoder, error) {
	return p.wotServers[ctxPath].GetEncoder(encoding)
}

func (p *Http) actionStartHandler(wotServer *WotServer, actionName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		encoder, err := wotServer.GetEncoder(encoder.ENCODING_JSON)

		if err != nil {
			sendJsonERR(w, err)
			return
		}

		wo := WotObject{}
		err = encoder.Unmarshal(r.Body, &wo)

		if err != nil {
			sendERR(w, encoder, err)
			return
		}

		actionID, slot := p.actionResults.CreateSlot()
		clients := async.NewFanOut()
		p.subscribers.CreateSubscription(actionID, clients)
		ph := NewProgressHandler(slot, clients)
		_, rc := wotServer.InvokeAction(actionName, wo, ph)

		if rc == WOT_OK {
			hrefs := links(websocketSubUrl(r, actionID), httpSubUrl(r, actionID))
			sendOK(w, encoder, hrefs)
		} else {
			sendERR(w, encoder, rc)
		}
	}
}

func (p *Http) actionTaskHandler(wotServer *WotServer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		encoder, err := wotServer.GetEncoder(encoder.ENCODING_JSON)

		if err != nil {
			sendJsonERR(w, err)
			return
		}

		vars := mux.Vars(r)
		taskid := vars["taskid"]
		slot, rc := p.actionResults.GetSlot(taskid)

		if rc {
			sendOK(w, encoder, slot.Load())
		} else {
			sendERR(w, encoder, rc)
		}
	}
}

func (p *Http) actionWSTaskHandler(wotServer *WotServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		taskid := vars["taskid"]
		p.wsHandler(wotServer, taskid, w, r)
	}
}

func (p *Http) wsHandler(wotServer *WotServer, handlerId string, w http.ResponseWriter, r *http.Request) {
	encoder, err := wotServer.GetEncoder(encoder.ENCODING_JSON)

	if err != nil {
		sendJsonERR(w, err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Error creating WebSocket at: ", err)
		return
	}

	client := make(chan interface{})
	clientID := p.subscribers.AddClient(handlerId, client)

	log.Println("Created internal subscriber handlerId: ", handlerId, " clientID: ", clientID)

	wsOpened := true
	for event := range client {
		if err = writeJSON(conn, encoder, event); err != nil && wsOpened {
			p.subscribers.RemoveClient(handlerId, clientID)
			log.Println("Removed internal subscriber handlerId: ", handlerId, " clientID: ", clientID)
			wsOpened = false
		}
	}
}

// CREDIT TO Gorilla websocket library
func writeJSON(wsc *websocket.Conn, codec encoder.Encoder, v interface{}) error {
	w, err := wsc.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	err1 := codec.Marshal(w, v)
	err2 := w.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (p *Http) eventSubscribeHandler(wotServer *WotServer, eventName string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		encoder, err := wotServer.GetEncoder(encoder.ENCODING_JSON)

		if err != nil {
			sendJsonERR(w, err)
			return
		}

		subscriptionID, _ := sec.UUID4()
		clients := async.NewFanOut()

		p.subscribers.CreateSubscription(subscriptionID, clients)
		wotServer.AddListener(eventName, p.eventHandler(subscriptionID, clients))

		hrefs := links(websocketSubUrl(r, subscriptionID))
		sendOK(w, encoder, hrefs)
	}
}

func (p *Http) eventWSClientHandler(wotServer *WotServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		subscriptionID := vars["subscriptionID"]
		p.wsHandler(wotServer, subscriptionID, w, r)
	}
}

func (p *Http) eventHandler(uuid string, clients *async.FanOut) *EventListener {
	el := &EventListener{
		ID: uuid,
		CB: func(event interface{}) {
			clients.Publish(event)
		},
	}

	return el
}

func links(links ...Link) *Links {
	ls := Links{
		Links: make([]Link, 0),
	}

	for _, l := range links {
		ls.Links = append(ls.Links, l)
	}

	return &ls
}

func httpSubUrl(r *http.Request, subresource string) Link {
	uri := removeTTslash(r.URL.RequestURI())

	if len(uri) == 0 {
		uri = str.Concat("/", removeTTslash(subresource))
	} else {
		uri = str.Concat("/", uri, "/", removeTTslash(subresource))
	}

	linkString := str.Concat("http://", r.Host, uri)

	return Link{
		Rel:  "rest",
		Href: linkString,
	}
}

func websocketSubUrl(r *http.Request, subresource string) Link {
	uri := removeTTslash(r.URL.RequestURI())

	if len(uri) == 0 {
		uri = str.Concat("/ws/", removeTTslash(subresource))
	} else {
		uri = str.Concat("/", uri, "/ws/", removeTTslash(subresource))
	}

	linkString := str.Concat("ws://", r.Host, uri)

	return Link{
		Rel:  "websocket",
		Href: linkString,
	}
}

func removeTTslash(str string) string {
	s := str

	if len(s) == 0 {
		return s
	}

	if s[len(s)-1:] == "/" {
		s = s[0 : len(s)-1]
	}

	if len(s) == 0 {
		return s
	}

	if s[0:1] == "/" {
		s = s[1:len(s)]
	}

	return s
}

func contextPath(ctxPath, element string) string {
	return str.Concat(ctxPath, "/", element)
}

func sendOK(w http.ResponseWriter, codec encoder.Encoder, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	codec.Marshal(w, payload)
}

func sendJsonOK(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(payload)
}

func sendJsonERR(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	switch payload.(type) {
	default:
		json.NewEncoder(w).Encode(payload)
	case error:
		json.NewEncoder(w).Encode(payload.(error).Error())
	}
}

func sendERR(w http.ResponseWriter, encoder encoder.Encoder, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	switch payload.(type) {
	default:
		encoder.Marshal(w, payload)
	case error:
		encoder.Marshal(w, payload.(error).Error())
	}
}

type Links struct {
	Links []Link `json:"links"`
}

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type route struct {
	name        string
	method      string
	pattern     string
	handlerFunc http.HandlerFunc
}

func (p *Http) addRoute(route *route) {
	p.router.
		Methods(route.method).
		Path(route.pattern).
		Name(route.name).
		Handler(route.handlerFunc)
}
