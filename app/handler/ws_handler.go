package handler

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type (
	WsConn struct {
		*websocket.Conn
	}

	WsItem struct {
		State   string `json:"state"`
		Message string `json:"message"`
	}

	WsPayload struct {
		Size    int    `json:"size"`
		Message string `json:"message"`
	}
)

var (
	clients      = make(map[WsConn]string)
	currentCount = 0
	currentSize  = 0
	wsChan       = make(chan WsPayload)
	mx           = sync.RWMutex{}
	upgrader     = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

func WebSocketHandler(w http.ResponseWriter, r *http.Request, workerPool chan struct{}) {
	<-workerPool

	defer func() {
		workerPool <- struct{}{}
	}()

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("[socket-aggregator] ws. connection upgrade error:", err)
		return
	}

	response := WsItem{
		State:   "processor",
		Message: "message received",
	}

	mx.Lock()
	conn := WsConn{Conn: ws}
	clients[conn] = ""

	go func() {
		if err = ws.WriteJSON(response); err != nil {
			fmt.Println("[socket-aggregator] write-json error:", err)
		}
	}()
	mx.Unlock()

	go wsListener(&conn)
	go wsChannelListener()
}

// HELPER METHODS

func wsListener(conn *WsConn) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var wsItem WsItem

	for {
		err := conn.ReadJSON(&wsItem)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				_ = conn.Close()
				fmt.Println("[socket-processor] unexpected close error:", err)
			}
			break
		}

		var payload WsPayload
		payload.Size = int(uintptr(len(wsItem.Message)))
		payload.Message = wsItem.Message
		wsChan <- payload

	}
}

func wsChannelListener() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case ws := <-wsChan:
			currentCount++
			currentSize += ws.Size
		case <-time.After(3 * time.Second):
			stdoutUpdater(currentCount, currentSize)
		case <-signalChan:
			fmt.Printf("\nterminating the service...\n")
			close(wsChan)
			os.Exit(0)
		}
	}
}

func stdoutUpdater(values ...interface{}) {
	totalCount := fmt.Sprintf("[total count.....(%d)]", values[0])

	kb := values[1].(int) / 1024
	totalSize := fmt.Sprintf("[total size.....(%d KB)]", kb)

	fmt.Printf("\r%s#%s", totalCount, totalSize)
}
