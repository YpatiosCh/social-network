package clients

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

var connIdCounter atomic.Int64 // used to give each client connection its own id

var (
	// ErrDuplicateAdd       = fmt.Errorf("client seemingly already connected, both old and new conn rejected")
	ErrTooManyConnections = fmt.Errorf("too many connections for client")
	ErrClientNotFound     = fmt.Errorf("couldn't find client")
	ErrConnectionNotFound = fmt.Errorf("couldn't find connection")
	ErrUnlucky            = fmt.Errorf("can't add connection cause client was just removed, unlucky data race, try again")
	ErrCantMarshalData    = fmt.Errorf("can't marshal data")
)

const CLIENT_CONNECTION_LIMIT = 10 // how many connections (websockets) each client can have

type client struct {
	connMu      sync.RWMutex
	connections map[int64]*ClientConn
	deleted     bool // flag that means the client was removed from clients map, used for races where goroutine tries to add connection to removed client
}

type ClientConn struct {
	Channel chan<- models.InnerPacket
	Ctx     context.Context
}

var clients sync.Map

// adds a connection under its client's entry so that other parts of the code can send data to this specific connection.
// If client entry is missing we create it, there's a limit for how many connections (websockets) each client can have.
func AddClientConnection(clientId int64, clientChannel chan<- models.InnerPacket, ctx context.Context) (int64, error) {
	// check if client already exists
	clientAny, ok := clients.Load(clientId)
	var loaded bool

	if !ok {
		// client entry not found
		// so we'll create both client and connection

		// create new connection
		newId := connIdCounter.Add(1)
		connMap := make(map[int64]*ClientConn)
		connMap[newId] = &ClientConn{clientChannel, ctx}

		// create new client
		theClient := client{sync.RWMutex{}, connMap, false}
		clientAny, loaded = clients.LoadOrStore(clientId, &theClient)
		if !loaded {
			return newId, nil
		}
		// reaching here means there was a race condition, and someone elses made the client before us, so we'll try to add this connection to the client again
	}

	// client entry found
	// so we'll add a new connection to the map

	c, ok := clientAny.(*client)
	if !ok {
		panic(1)
	}

	c.connMu.Lock()
	defer c.connMu.Unlock()
	// check if the client had 0 connections and it was just deleted
	if c.deleted {
		return 0, ErrUnlucky
	}

	// check for too many connections
	if len(c.connections) >= CLIENT_CONNECTION_LIMIT {
		fmt.Println("connections", len(c.connections))
		return 0, ErrTooManyConnections
	}
	// add new connection
	newId := connIdCounter.Add(1)
	c.connections[newId] = &ClientConn{clientChannel, ctx}

	return newId, nil
}

// removes client connection entry
// if all connections are removed, client entry also gets removed
func RemoveClientConnection(clientId int64, connectionId int64) error {
	clientAny, ok := clients.Load(clientId)
	if !ok {
		return ErrClientNotFound
	}
	c, ok := clientAny.(*client)
	if !ok {
		panic(1)
	}

	c.connMu.Lock()
	defer c.connMu.Unlock()

	delete(c.connections, connectionId)
	// if no more connections, remove client from clients
	if len(c.connections) == 0 {
		c.deleted = true
		clients.Delete(clientId)
	}

	return nil
}

func GetClientConnection(clientId, connectionId int64) (*ClientConn, error) {
	anyClient, ok := clients.Load(clientId)
	if !ok {
		return nil, ErrClientNotFound
	}
	c, ok := anyClient.(*client)
	if !ok {
		panic(1)
	}

	c.connMu.RLock()
	defer c.connMu.RUnlock()
	conn, ok := c.connections[connectionId]
	if !ok {
		return nil, ErrConnectionNotFound
	}
	return conn, nil
}

// sends data to all connections of specific user. Either put your the data in json form in data arg, or the raw struct into anyData
//
// # If you have the json already encoded put it in jsonData
//
// # If you have just a struct, then leave data emtpy and put your struct to anyData
func SendDataToAllClientConns(clientId int64, anyData any, highPriority bool) error {
	anyClient, ok := clients.Load(clientId)
	if !ok {
		return ErrClientNotFound
	}
	c, ok := anyClient.(*client)
	if !ok {
		panic(1)
	}

	c.connMu.RLock()
	if len(c.connections) == 0 {
		c.connMu.RUnlock()
		return nil
	}
	connectionCopy := make([]*ClientConn, 0, len(c.connections))
	for _, conn := range c.connections {
		connectionCopy = append(connectionCopy, conn)
	}
	c.connMu.RUnlock()

	for _, conn := range connectionCopy {
		select {
		case conn.Channel <- models.InnerPacket{Payload: anyData, HighPriority: highPriority}:
		case <-conn.Ctx.Done():
		case <-time.After(time.Second):
		}
	}

	return nil
}

// Sends data to LITERALLY EVERYONE!!
//
// # If you have the json already encoded put it in jsonData
//
// # If you have just a struct, then leave data emtpy and put your struct to anyData
func SendToAllClients(anyData any, highPriority bool) error {

	failCounter := 0

	clients.Range(func(anyClientId, anyClient any) bool {
		clientId, ok := anyClientId.(int64)
		if !ok {
			panic(1)
		}

		err := SendDataToAllClientConns(clientId, anyData, highPriority)
		if err != nil {
			fmt.Println("Failed to send Post to client:", err)
			failCounter++
		}
		return true
	})

	if failCounter > 0 {
		fmt.Printf("Sending data to all clients: failed to send to %d clients\n", failCounter)
	}

	return nil
}

func UpdateUserStatuses(dbResp *models.UsersFeedDbResponse) {
	for i := range dbResp.Feed {
		userId := dbResp.Feed[i].User.Id

		if _, ok := clients.Load(userId); ok {
			dbResp.Feed[i].User.Status = "online"
		} else {
			dbResp.Feed[i].User.Status = "offline"
		}
	}
}

func UpdateUserStatus(feed *models.Feed) {
	if _, ok := clients.Load(feed.User.Id); ok {
		feed.User.Status = "online"
	} else {
		feed.User.Status = "offline"
	}
}

func CheckMulttipleConnections(userId int64) bool {

	clientInterface, exists := clients.Load(userId)
	if !exists {
		fmt.Println("Client not found")
		return false
	}

	// Type assert the value to the correct type
	client, ok := clientInterface.(*client)
	if !ok {
		fmt.Println("Failed to assert client type")
		return false
	}

	client.connMu.RLock()
	defer client.connMu.RUnlock()

	// Check if the client has more than one connection
	if len(client.connections) > 1 {
		return true
	}

	fmt.Println("Client has one or fewer connections")
	return false

}
