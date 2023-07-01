package main

import (
	"fmt"

	"github.com/jtribble/fly-io-dist-sys/pkg/maelstrom"
	"github.com/jtribble/fly-io-dist-sys/pkg/ptr"
)

type GenerateHandler struct {
	nextId int
}

func (h *GenerateHandler) HandlesMessageType(msgType string) bool {
	return msgType == "generate"
}

func (h *GenerateHandler) HandleMessage(node *maelstrom.Node, msg *maelstrom.Message) {
	node.SendMessage(&maelstrom.Message{
		Body: maelstrom.MessageBody{
			Type: "generate_ok",
			Id:   ptr.ToString(fmt.Sprintf("%s-%d", node.Id(), h.nextId)),
		},
	}, msg, nil)
	h.nextId += 1
}
