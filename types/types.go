package types

type Message struct {
	Action  string `json:"action"`
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
}

type Absent struct {
	Action   string `json:"action"`
	SenderID string `json:"sender_id"`
}

type StatusUpdate struct {
	Action string `json:"action"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type CheckStatus struct {
	Action string `json:"action"`
	Chash  string `json:"c_hash"`
}

type StatusResponse struct {
	Action   string `json:"action"`
	Chash    string `json:"c_hash"`
	LastSeen string `json:"last_seen"`
}

type Header interface {
	Type() string
}

// Message struct represents the structure of the JSON message

func (n Message) Type() string {
	return n.Action
}

func (n Absent) Type() string {
	return n.Action
}
func (n StatusUpdate) Type() string {
	return n.Action
}
func (n CheckStatus) Type() string {
	return n.Action
}
func (n StatusResponse) Type() string {
	return n.Action
}
