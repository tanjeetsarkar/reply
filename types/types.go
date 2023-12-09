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

type Token struct {
	user  string
	token string
}
