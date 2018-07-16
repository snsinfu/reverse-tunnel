package service

// BinderAcceptMessage is a model of a JSON message sent via websocket on client
// connection.
type BinderAcceptMessage struct {
	Event       string `json:"event"`
	SessionID   string `json:"session_id"`
	PeerAddress string `json:"peer_address"`
}
