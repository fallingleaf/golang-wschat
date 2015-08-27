package pubsub

type custom interface{}
type map_data map[string] custom

/*
 * Type chat: data {topic_id, message}
 * Type method: data {api, params}
 */
type Message struct {
    Type int `json:"type"`
    Data map_data `json:"data"`
}

type Event struct {
    etype int
    source string
    data custom
}

const (
    // Available event types
    e_status = 1
    e_message = 2
    e_sync = 3
)

const (
    // Available status types
    status_offline = 0
    status_online = 1
    status_away = 2
    status_invi = 3
    status_busy = 4
)

const (
    // Available message types
    type_chat = 0
    type_method = 1
)

// Register handlers as map
var methods = map[string] func(map_data) error {
    "auth": auth,
    "subscribe": subscribe,
}
