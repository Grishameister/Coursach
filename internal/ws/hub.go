package ws

import (
	"encoding/json"
	"github.com/Grishameister/Coursach/internal/domain"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var writeWait = 1 * time.Second

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	statusMessages chan domain.StatusChannel
	Connections    map[uuid.UUID]*websocket.Conn
	mu             *sync.Mutex
}

func NewHub(statusMessages chan domain.StatusChannel) *Hub {
	return &Hub{
		statusMessages: statusMessages,
		Connections:    make(map[uuid.UUID]*websocket.Conn),
		mu:             &sync.Mutex{},
	}
}

func (h *Hub) addConnection(ws *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	id, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}
	h.Connections[id] = ws
}

func (h *Hub) ServeWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	h.addConnection(ws)
}

func (h *Hub) deleteConn(id uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if ws, ok := h.Connections[id]; ok {
		defer delete(h.Connections, id)
		if err := ws.Close(); err != nil {
			log.Println(err)
		}
	}
}

func (h *Hub) WriteMessages() {
	for m := range h.statusMessages {
		rawMessage, err := json.Marshal(&m)
		if err != nil {
			log.Println(err)
			continue
		}

		h.mu.Lock()
		for id, ws := range h.Connections {
			log.Println(id)
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteMessage(websocket.TextMessage, rawMessage)
			if err != nil && err != websocket.ErrCloseSent {
				log.Println(err)
				h.deleteConn(id)
			}
		}
		h.mu.Unlock()
	}
}
