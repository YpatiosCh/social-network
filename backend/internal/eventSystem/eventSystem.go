package eventSystem

import (
	"errors"
	"fmt"
	"hash/maphash"
	"maps"
	"sync"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/clients"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

var (
	// ErrClientConnectionNotFound = errors.New("client connection not found")
	ErrEventNotFound   = errors.New("couldn't find event")
	ErrCantMarshalData = errors.New("can't marshal data")
	ErrDoublicateSub   = errors.New("doublicate subscription")
	ErrUnlucky         = errors.New("unlucky, variable was 'deleted' between capturing and interacting, try again")
)

// this will be used to split the event map into X number of smaller maps
const shardCount = 100

// holds the shards that contain 'event name' - 'subscribers' slice pairs for each event
var eventShards = []*sync.Map{}

// holds clients connectionId int64 as key, and the events its subscribed to as values. Used for purging events of specific connection or all client connections
var clientConnectionShards = []*sync.Map{}

// holds a map of the subscribers that are listening to a particular event
type subscribers struct {
	mu      sync.RWMutex
	subMap  map[int64]*clients.ClientConn
	deleted bool
}

// holds a map of connections and what events they are subscribed to
type subscriptions struct {
	mu      sync.RWMutex
	subMap  map[string]struct{}
	deleted bool
}

func init() {
	for range shardCount {
		eventShards = append(eventShards, &sync.Map{})
		clientConnectionShards = append(clientConnectionShards, &sync.Map{})
	}

}

func AddSubscription(clientId, connectionId int64, eventKey string) error {
	fmt.Println(connectionId, "subbed to", eventKey)
	// add client connection to event entry
	err := addConnectionToEventMap(clientId, connectionId, eventKey)
	if err != nil {
		return err
	}

	// add subscription to connectionId, used for finding all subscriptions belonging to a connection or client
	err = addEventToConnectionMap(connectionId, eventKey)
	if err != nil {
		return err
	}

	PrintSubscriptions()
	return nil
}

func addConnectionToEventMap(clientId, connectionId int64, eventKey string) error {
	events := getEventShard(eventKey)

	// check if event entry already exists
	value, ok := events.Load(eventKey)
	if !ok {
		// no event entry exists, so we'll add a new one
		s := &subscribers{sync.RWMutex{}, make(map[int64]*clients.ClientConn), false}
		value, _ = events.LoadOrStore(eventKey, s)
	}

	subs, ok := value.(*subscribers)
	if !ok {
		panic(1)
	}

	subs.mu.Lock()
	defer subs.mu.Unlock()

	if subs.deleted {
		return ErrUnlucky
	}

	_, ok = subs.subMap[connectionId]
	if ok {
		return ErrDoublicateSub
	}

	clientConnection, err := clients.GetClientConnection(clientId, connectionId)
	if err != nil {
		return fmt.Errorf("can't add subscription: %w", err)
	}
	subs.subMap[connectionId] = clientConnection

	return nil
}

func addEventToConnectionMap(connectionId int64, eventKey string) error {
	connectionEvents := getConnectionShard(connectionId)
	anyConnEvents, ok := connectionEvents.Load(connectionId)
	if !ok {
		m := subscriptions{sync.RWMutex{}, make(map[string]struct{}), false}
		anyConnEvents, _ = connectionEvents.LoadOrStore(connectionId, &m)
	}
	connectionSubscriptions, ok := anyConnEvents.(*subscriptions)
	if !ok {
		panic(1)
	}

	connectionSubscriptions.mu.Lock()
	defer connectionSubscriptions.mu.Unlock()

	if connectionSubscriptions.deleted {
		return ErrUnlucky
	}

	connectionSubscriptions.subMap[eventKey] = struct{}{}

	return nil
}

func RemoveSubscription(connectionId int64, eventKey string) error {
	fmt.Println(connectionId, "unsubbed from", eventKey)
	// remove connection from events
	eventShard := getEventShard(eventKey)
	anySubscribers, ok := eventShard.Load(eventKey)
	if !ok {
		return ErrEventNotFound
	}
	subscribers, ok := anySubscribers.(*subscribers)
	if !ok {
		panic(1)
	}

	subscribers.mu.Lock()

	delete(subscribers.subMap, connectionId)
	if len(subscribers.subMap) == 0 {
		subscribers.deleted = true
		eventShard.Delete(eventKey)
	}

	subscribers.mu.Unlock()

	// remove event from connections
	connEventsShard := getConnectionShard(connectionId)
	anyConnEvents, ok := connEventsShard.Load(connectionId)
	if !ok {
		return ErrEventNotFound
	}
	subscriptions, ok := anyConnEvents.(*subscriptions)
	if !ok {
		panic(1)
	}
	subscriptions.mu.Lock()
	delete(subscriptions.subMap, eventKey)
	if len(subscriptions.subMap) == 0 {
		subscriptions.deleted = true
		connEventsShard.Delete(connectionId)
	}
	subscriptions.mu.Unlock()

	PrintSubscriptions()
	return nil
}

// Sends data to all connections that are subscribed to event
//
// # If you have the json already encoded put it in jsonData
//
// # If you have just a struct, then leave data emtpy and put your struct to anyData
func SendEvent(eventKey string, anyData any, highPriority bool) error {

	events := getEventShard(eventKey)

	anySubscribers, ok := events.Load(eventKey)
	if !ok {
		return ErrEventNotFound
	}
	subscribers, ok := anySubscribers.(*subscribers)
	if !ok {
		panic(1)
	}

	subscribers.mu.RLock()
	channelSlice := make([]*clients.ClientConn, 0, len(subscribers.subMap))
	for _, clientConnection := range subscribers.subMap {
		channelSlice = append(channelSlice, clientConnection)
	}
	subscribers.mu.RUnlock()

	for _, clientConnection := range channelSlice {
		select {
		case clientConnection.Channel <- models.InnerPacket{Payload: anyData, HighPriority: highPriority}:
		case <-clientConnection.Ctx.Done():
			// TODO think what I need to do here
		}
	}

	return nil
}

// seed that will be used for hashing function, all goroutines will use the same one
var seed = maphash.MakeSeed()

// sync pool of hashing objects, since each goroutine will need its own one
var hasherPool = sync.Pool{New: func() any {
	var h maphash.Hash
	h.SetSeed(seed)
	return &h
}}

// given a key, it will provide the appropriate shard index
func key2Shard(key string) uint64 {
	h, ok := hasherPool.Get().(*maphash.Hash)
	if !ok {
		panic(1)
	}
	h.Reset()
	h.SetSeed(seed)

	_, err := h.WriteString(key)
	if err != nil {
		panic(err)
	}
	finalValue := h.Sum64() % shardCount

	hasherPool.Put(h)

	return finalValue
}

// returns the shard that belong to the key
func getEventShard(eventKey string) *sync.Map {
	return eventShards[int(key2Shard(eventKey))]
}

// returns the shard that belong to the key
func getConnectionShard(connectionId int64) *sync.Map {
	return clientConnectionShards[int(key2Shard(fmt.Sprint(connectionId)))]
}

func validateEvent(clientId int64, eventKey string) error {
	return nil
}

func PurgeSubscriptions(connectionId int64) error {
	connEventsShard := getConnectionShard(connectionId)
	anyConnEvents, ok := connEventsShard.Load(connectionId)
	if !ok {
		return ErrEventNotFound
	}
	subscriptions, ok := anyConnEvents.(*subscriptions)
	if !ok {
		panic(1)
	}
	subscriptions.mu.RLock()
	keys := maps.Keys(subscriptions.subMap)
	subscriptions.mu.RUnlock()

	for key := range keys {
		err := RemoveSubscription(connectionId, key)
		if err != nil {
			fmt.Printf("warning: failed to remove sub: %s for clientid %d", key, connectionId)
		}
	}

	return nil
}

func CreateDmEventKey(idA, idB int64) string {
	var eventKey string
	if idA < idB {
		eventKey = fmt.Sprintf("dm:%d:%d", idA, idB)
	} else {
		eventKey = fmt.Sprintf("dm:%d:%d", idB, idA)
	}
	return eventKey
}

func PrintSubscriptions() {
	return
	fmt.Println("==== Event Subscriptions ====")
	for i, shard := range eventShards {
		shard.Range(func(key, value any) bool {
			eventKey, ok := key.(string)
			if !ok {
				return true
			}
			subs, ok := value.(*subscribers)
			if !ok {
				return true
			}
			subs.mu.RLock()
			ids := make([]int64, 0, len(subs.subMap))
			for connId := range subs.subMap {
				ids = append(ids, connId)
			}
			subs.mu.RUnlock()
			fmt.Printf("Shard %d | Event: %s | Connections: %v\n", i, eventKey, ids)
			return true
		})
	}
	fmt.Println("==== Connection Subscriptions ====")
	for i, shard := range clientConnectionShards {
		shard.Range(func(key, value any) bool {
			connId, ok := key.(int64)
			if !ok {
				return true
			}
			subs, ok := value.(*subscriptions)
			if !ok {
				return true
			}
			subs.mu.RLock()
			events := make([]string, 0, len(subs.subMap))
			for eventKey := range subs.subMap {
				events = append(events, eventKey)
			}
			subs.mu.RUnlock()
			fmt.Printf("Shard %d | Connection: %d | Events: %v\n", i, connId, events)
			return true
		})
	}
}
