package msg

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Message struct {
	ID string
}

func New() Message {
	id := uuid.New().String()
	return Message{ID: id}
}

func (m Message) ToJSON() ([]byte, error) {
	bs, err := json.Marshal(&m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON; %w", err)
	}
	return bs, nil
}

func FromJSON(bs []byte) (Message, error) {
	m := Message{}
	if err := json.Unmarshal(bs, &m); err != nil {
		return m, fmt.Errorf("failed to unmarshal JSON; %w", err)
	}
	return m, nil
}
