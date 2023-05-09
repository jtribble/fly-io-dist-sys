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
		node.Id = *msg.Body.NodeId
		node.Peers = slices.Without(*msg.Body.NodeIds, node.Id)
		node.QueueReply(&Message{
			Body: MessageBody{
				Type: "init_ok",
			},
		}, msg)
	}
	h.initialized = true
}
