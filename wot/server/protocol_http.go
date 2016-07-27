package server

import (
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

// FIXMEs:

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

	updateThingDescription(ctxPath, td)
}

func updateThingDescription(ctxPath string, td *model.ThingDescription) {
	td.Uris = append(td.Uris, str.Concat(hostname, ctxPath))
	td.Encodings = encoder.Registry.Registered()
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
		method:  "GET",
		pattern: "/",
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
			ls := links()

			for path := range p.wotServers {
				ls.Links = append(ls.Links, httpSubUrl(r, path))
			}

			sendOK(w, r, ls)
		},
	})
}

func (p *Http) registerDeviceRoot(ctxPath string, td *model.ThingDescription) {
	p.addRoute(&route{
		method:  "GET",
		pattern: contextPath(ctxPath, ""),
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
			hrefs := links(httpSubUrl(r, "description"))

			sendOK(w, r, hrefs)
		},
	})
}

func (p *Http) registerDeviceDescriptor(ctxPath string, td *model.ThingDescription) {
	p.addRoute(&route{
		method:  "GET",
		pattern: contextPath(ctxPath, "description"),
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
			sendOK(w, r, td)
		},
	})
}

func (p *Http) registerProperties(ctxPath string, td *model.ThingDescription) {
	for _, prop := range td.Properties {
		p.addRoute(&route{
			method:      "GET",
			pattern:     contextPath(ctxPath, prop.Hrefs[0]),
			handlerFunc: p.propertyGetHandler(ctxPath, &prop),
		})

		if prop.Writable {
			p.addRoute(&route{
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
			method:      "POST",
			pattern:     contextPath(ctxPath, action.Hrefs[0]),
			handlerFunc: p.actionStartHandler(p.wotServers[ctxPath], action.Name),
		})

		p.addRoute(&route{
			method:      "GET",
			pattern:     contextPath(ctxPath, str.Concat(action.Hrefs[0], "/{taskid}")),
			handlerFunc: p.actionTaskHandler(p.wotServers[ctxPath]),
		})

		p.addRoute(&route{
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
			method:      "POST",
			pattern:     contextPath(ctxPath, event.Hrefs[0]),
			handlerFunc: p.eventSubscribeHandler(p.wotServers[ctxPath], event.Name),
		})

		p.addRoute(&route{
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
		value := p.wotServers[ctxPath].GetProperty(prop.Name)

		data := value.Get()

		switch data.(type) {
		case Status:
			if data.(Status) != WOT_OK {
				sendERR(w, r, data)
			}
		case error:
			sendERR(w, r, data)
		default:
			sendOK(w, r, data)
		}
	}
}

func (p *Http) propertySetHandler(ctxPath string, prop *model.Property) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		wo := WotObject{}
		err := readBody(r, &wo)

		if err != nil {
			sendPlainERR(w, err)
			return
		}

		value := p.wotServers[ctxPath].SetProperty(prop.Name, wo.GetValue())
		data := value.Get()

		switch data.(type) {
		case Status:
			if data.(Status) != WOT_OK {
				sendERR(w, r, data)
			}
		case error:
			sendERR(w, r, data)
		}
	}
}

func (p *Http) actionStartHandler(wotServer *WotServer, actionName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		wo := WotObject{}
		err := readBody(r, &wo)

		if err != nil {
			sendPlainERR(w, err)
			return
		}

		actionID, slot := p.actionResults.CreateSlot()
		clients := async.NewFanOut()
		p.subscribers.CreateSubscription(actionID, clients)
		ph := NewProgressHandler(slot, clients)
		value := wotServer.InvokeAction(actionName, wo, ph)
		data := value.Get()

		switch data.(type) {
		case Status:
			if data.(Status) != WOT_OK {
				sendERR(w, r, data)

			}
		case error:
			sendERR(w, r, data)
		default:
			hrefs := links(websocketSubUrl(r, actionID), httpSubUrl(r, actionID))
			sendOK(w, r, hrefs)
		}
	}
}

func (p *Http) actionTaskHandler(wotServer *WotServer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		taskid := vars["taskid"]
		slot, rc := p.actionResults.GetSlot(taskid)

		if rc {
			sendOK(w, r, slot.Load())
		} else {
			sendERR(w, r, rc)
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
		if err = writeJSON(conn, r, event); err != nil && wsOpened {
			p.subscribers.RemoveClient(handlerId, clientID)
			log.Println("Removed internal subscriber handlerId: ", handlerId, " clientID: ", clientID)
			wsOpened = false
		}
	}
}

// CREDIT TO Gorilla websocket library
func writeJSON(wsc *websocket.Conn, r *http.Request, v interface{}) error {
	w, err := wsc.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	encoder, err := encoder.Registry.Get("JSON")
	if err != nil {
		w.Write([]byte("Unsupported Encoding: JSON"))
		return err
	}

	err1 := encoder.Marshal(w, v)
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
		subscriptionID, _ := sec.UUID4()
		clients := async.NewFanOut()

		p.subscribers.CreateSubscription(subscriptionID, clients)
		wotServer.AddListener(eventName, p.eventHandler(subscriptionID, clients))

		hrefs := links(websocketSubUrl(r, subscriptionID))
		sendOK(w, r, hrefs)
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

func readBody(r *http.Request, t interface{}) error {
	encoder, err := encoder.Registry.Get("JSON")

	if err != nil {
		return err
	}

	err = encoder.Unmarshal(r.Body, t)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func sendOK(w http.ResponseWriter, r *http.Request, payload interface{}) {
	encoder, err := encoder.Registry.Get("JSON")

	if err != nil {
		sendPlainERR(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder.Marshal(w, payload)
}

func sendERR(w http.ResponseWriter, r *http.Request, payload interface{}) {
	encoder, err := encoder.Registry.Get("JSON")

	if err != nil {
		sendPlainERR(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	switch payload.(type) {
	default:
		encoder.Marshal(w, payload)
	case error:
		encoder.Marshal(w, payload.(error).Error())
	}
}

func sendPlainERR(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)

	w.Write([]byte(err.Error()))
}

type Links struct {
	Links []Link `json:"links"`
}

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type route struct {
	method      string
	pattern     string
	handlerFunc http.HandlerFunc
}

func (p *Http) addRoute(route *route) {
	p.router.
		Methods(route.method).
		Path(route.pattern).
		Name(route.pattern).
		Handler(route.handlerFunc)
}
