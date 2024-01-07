package validation

import (
	"encoding/json"
	"fmt"

	"github.com/reply/types"
)

func ValidateAction(jsonData []byte) (types.Header, error) {

	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Println("Invalid data received")
		return nil, fmt.Errorf("invalid data received")
	}

	action, ok := data["action"].(string)
	if !ok {
		return nil, fmt.Errorf("no action received")
	}

	switch action {
	case "TEXT_MESSAGE":
		var message types.Message
		err := json.Unmarshal(jsonData, &message)
		if err != nil {
			return nil, fmt.Errorf("invalid message data received")
		}
		return message, nil
	case "ABSENT":
		var absent types.Absent
		err := json.Unmarshal(jsonData, &absent)
		if err != nil {
			return nil, fmt.Errorf("invalid absent data received")
		}
		return absent, nil
	case "USER_JOIN":
		var status_update types.StatusUpdate
		err := json.Unmarshal(jsonData, &status_update)
		if err != nil {
			return nil, fmt.Errorf("invalid absent data received")
		}
		return status_update, nil
	default:
		return nil, fmt.Errorf("invalid default data received")
	}
}
