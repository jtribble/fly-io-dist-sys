package maelstrom

import (
	"github.com/jtribble/fly-io-dist-sys/pkg/slices"
)

type InitHandler struct {
	initialized bool
}

func (h *InitHandler) HandlesMessageType(msgType string) bool {
	return msgType == "init"
}

func (h *InitHandler) HandleMessage(node *Node, msg *Message) {
	if !h.initialized {
		node.id = *msg.Body.NodeId
		node.peers = slices.Without(*msg.Body.NodeIds, node.id)
		node.SendMessage(&Message{
			Body: MessageBody{
				Type: "init_ok",
			},
		}, msg, nil)
	}
	h.initialized = true
}
