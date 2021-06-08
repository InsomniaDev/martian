package zwave

var (
	// EventTopic is the main topic used for communication
	EventTopic = "zw/_EVENTS/ZWAVE_GATEWAY-server/value_changed"
)

// Zwave main struct
type Zwave struct {
	Messages []string
}