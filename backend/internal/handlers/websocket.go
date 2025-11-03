package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/clients"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/eventSystem"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/security"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/utils"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

// FOUND FROM DOCUMENTATION:

// The message types are defined in RFC 6455, section 11.8.

// TextMessage denotes a text data message. The text message payload is
// interpreted as UTF-8 encoded text data.
// TextMessage = 1

// BinaryMessage denotes a binary data message.
// BinaryMessage = 2

// CloseMessage denotes a close control message. The optional message
// payload contains a numeric code and text. Use the FormatCloseMessage
// function to format a close message payload.
// CloseMessage = 8

// PingMessage denotes a ping control message. The optional message payload
// is UTF-8 encoded text.
// PingMessage = 9

// PongMessage denotes a pong control message. The optional message payload
// is UTF-8 encoded text.
// PongMessage = 10

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Handlers) startWebsocket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("start websocket handler called")

		websocketConn, cancelConn, err := upgradeConnection(w, r)
		if err != nil {
			println(err.Error())
			return
		}
		defer cancelConn()

		claims, ok := utils.GetValue[security.Claims](r, utils.ClaimsKey)
		if !ok {
			panic(1)
		}

		clientId := claims.UserId
		fmt.Println("Client id:", clientId)

		// context for cancelling send goroutine when connection closes, and for other goroutines to not send anything to channel that is no longer needed
		ctx, cancelContext := context.WithCancel(context.Background())

		//adding client so that this client can be tracked and receive data
		channel, connectionId, err := addClientConnection(clientId, ctx)
		if err != nil {
			if errors.Is(err, clients.ErrTooManyConnections) {
				sendErrorToWS(websocketConn, "excessive_connections")
			}
			fmt.Println("error with adding client:", err.Error())
			cancelContext()
			return
		}

		// wait group to wait for sender goroutine to stop
		var wg sync.WaitGroup

		wg.Add(1)
		go websocketSender(connectionId, clientId, channel, websocketConn, ctx, &wg)

		websocketListener(websocketConn, clientId, connectionId)

		// CLOSING PRECEDURE

		// remove client so that no other goroutine attempts to send here
		clients.RemoveClientConnection(clientId, connectionId)

		//adding message about online status here, may need to move
		sendUserStatusUpdate(clientId, "offline")

		// remove all subscriptions associated with connectionId
		eventSystem.PurgeSubscriptions(connectionId)

		// close the context to stop sender goroutine and other goroutines from sending
		cancelContext()

		// wait for sender goroutine to stop
		wg.Wait()

		fmt.Println("ws handler closing")
	}
}

func sendErrorToWS(websocketConn *websocket.Conn, payload string) {
	errorMessage := []any{models.WSMessage{
		Type:    payload,
		Payload: "",
	}}

	bundledMessage, err := json.Marshal(errorMessage)
	if err != nil {
		fmt.Println("this isn't supposed to happen")
		panic(1)
	}

	err = websocketConn.WriteMessage(websocket.TextMessage, bundledMessage)
	if err != nil {
		fmt.Println("failed to inform user that they have too many tabs open")
	}
}

// routine that reads data coming from this client connection
func websocketListener(websocketConn *websocket.Conn, clientId int64, connectionId int64) {
	for {
		_, msg, err := websocketConn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				break
			}
			log.Printf("clientId: %d, Start websocket error: unexpected read error: %v\n", clientId, err)
			return
		}

		messageString := string(msg)
		if len(messageString) < 2 {
			fmt.Println("invalid message received:", messageString)
			return
		}

		action := messageString[0]
		payload := messageString[1:]

		var eventKey string
		switch {
		case strings.HasPrefix(payload, "dm:"):
			parts := strings.Split(payload, ":")
			otherPersonId, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil || len(parts) > 2 {
				fmt.Println("bad dm event key sent")
			} else {
				eventKey = eventSystem.CreateDmEventKey(clientId, otherPersonId)
			}
		default:
			eventKey = payload
		}

		if eventKey == "" {
			return
		}

		switch action {
		case '+':
			if err := eventSystem.AddSubscription(clientId, connectionId, eventKey); err != nil {
				fmt.Println("failed to add subscription:", err.Error())
			}
		case '-':
			eventSystem.RemoveSubscription(connectionId, eventKey)
		default:
			fmt.Println("unknown type of message received from websocket:", messageString)
		}

	}
}

func addClientConnection(clientId int64, ctx context.Context) (chan models.InnerPacket, int64, error) {

	// this connectionChannel will be used by various parts of the program to send data to this client
	connectionChannel := make(chan models.InnerPacket, 10)

	var connectionId int64 = 0
	var err error

	retry := 0
	keepTrying := true
	for keepTrying {
		connectionId, err = clients.AddClientConnection(clientId, connectionChannel, ctx)
		if err != nil {
			if errors.Is(err, clients.ErrUnlucky) {
				if retry > 3 {
					fmt.Println("too many retries! something is very wrong!")
					panic(1)
				}
				time.Sleep(time.Millisecond * 10)
				retry++
				continue

			} else {
				return nil, 0, err
			}
		}

		keepTrying = false
	}

	fmt.Printf("WebSocket registered for client %d with connectionId %d\n", clientId, connectionId)

	sendUserStatusUpdate(clientId, "online")

	return connectionChannel, connectionId, nil
}

func upgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, func(), error) {
	websocketConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "failed websocket upgrade")
		return nil, func() {}, fmt.Errorf("failed to upgrade connection to websocket: %w", err)
	}
	deferMe := func() {
		err := websocketConn.Close()
		if err != nil {
			fmt.Println("failed to close websocket connection")
		}
	}
	return websocketConn, deferMe, nil
}

// Goroutine that sends data to this connection
func websocketSender(connectionid, clientId int64, channel <-chan models.InnerPacket, conn *websocket.Conn, ctx context.Context, wg *sync.WaitGroup) {
	var timer *time.Timer
	defer wg.Done()
outerLoop:
	for {
		select {
		case message := <-channel:

			slice := []any{message.Payload}

			if message.HighPriority {
				timer = time.NewTimer(time.Millisecond)
			} else {
				timer = time.NewTimer(time.Millisecond * 500)
			}

		bundleLoop:
			for {
				select {
				case message := <-channel:
					slice = append(slice, message.Payload)
					if message.HighPriority {
						timer.Reset(time.Millisecond)
					}
				case <-timer.C:
					bundledMessage, err := json.Marshal(slice)
					if err != nil {
						fmt.Println("AAAAA! A BAD JSON WAS GIVEN! The double encode failed! In production this situation will probably be ignored")
						panic(1)
						// TODO attempt individual sends? and drop w/e fails with a warning message
					}

					fmt.Println("sending", string(bundledMessage), "to", connectionid)
					err = conn.WriteMessage(websocket.TextMessage, bundledMessage)
					if err != nil {
						fmt.Printf("connectionid: %d, error on send message: clientid:%d err:%s \n", connectionid, clientId, err.Error())
						break outerLoop
					}
					break bundleLoop
				}
			}

		case <-ctx.Done():
			if timer != nil {
				if !timer.Stop() {
					<-timer.C
				}
			}
			break outerLoop

		}
	}
	// TODO look into errors that can only happen when writing, this should cancel the entire connection
}

func sendUserStatusUpdate(clientId int64, status string) {

	newStatus := models.UserUpdateWS{
		Id:     clientId,
		Status: status,
	}

	message := models.WSMessage{
		Type:    "user_update",
		Payload: newStatus,
	}

	err := eventSystem.SendEvent(fmt.Sprintf("os:%d", clientId), message, false)
	if err != nil {
		fmt.Println("Online status: Failed to send event:", err.Error())
	}

}
