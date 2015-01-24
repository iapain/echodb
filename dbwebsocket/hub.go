package dbwebsocket

type hub struct {
  // Registered connections.
  connections map[*connection]bool
  // Inbound messages from the connections.
  broadcast chan []byte
  // Register requests from the connections.
  register chan *connection
  // Unregister requests from connections.
  unregister chan *connection

  name string
}

var hubs = make(map[string]hub)

func FetchOrInitHub(name string) hub {
  if _, ok := hubs[name]; ok {
    return hubs[name]
  } else {
    newhub := hub{
      broadcast:   make(chan []byte),
      register:    make(chan *connection),
      unregister:  make(chan *connection),
      connections: make(map[*connection]bool),
    }
    hubs[name] = newhub
    newhub.run()
    return newhub
  }
}

func (h *hub) run() {
  for {
    select {
    case c := <-h.register:
      h.connections[c] = true
    case c := <-h.unregister:
      if _, ok := h.connections[c]; ok {
        delete(h.connections, c)
        close(c.send)
      }
    case m := <-h.broadcast:
      for c := range h.connections {
        select {
          case c.send <- m:
          default:
            delete(h.connections, c)
            close(c.send)
        }
      }
    }
  }
}
