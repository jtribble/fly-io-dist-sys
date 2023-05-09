package maelstrom

import (
	"encoding/json"
)

type Message struct {
	Src  string      `json:"src"`
	Dest string      `json:"dest"`
	Body MessageBody `json:"body"`
}

func From(line []byte) (*Message, error) {
	m := Message{}
	if err := json.Unmarshal(line, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

type MessageBody struct {
	Type string `json:"type"`

	MsgId     *int      `json:"msg_id,omitempty"`
	InReplyTo *int      `json:"in_reply_to,omitempty"`
	NodeId    *string   `json:"node_id,omitempty"`
	NodeIds   *[]string `json:"node_ids,omitempty"`
	Code      *int      `json:"code,omitempty"`
	Echo      *string   `json:"echo,omitempty"`
}
