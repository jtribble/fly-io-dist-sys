package main

import (
	"github.com/jtribble/fly-io-dist-sys/pkg/maelstrom"
	"github.com/jtribble/fly-io-dist-sys/pkg/maps"
)

type BroadcastHandler struct {
	messages map[int]bool
	topology map[string][]string
}

func NewBroadcastHandler() *BroadcastHandler {
	return &BroadcastHandler{
		messages: make(map[int]bool),
		topology: make(map[string][]string),
	}
}

func (h *BroadcastHandler) HandlesMessageType(msgType string) bool {
	return msgType == "broadcast" || msgType == "read" || msgType == "topology"
}

func (h *BroadcastHandler) HandleMessage(node *maelstrom.Node, msg *maelstrom.Message) {
	switch msg.Body.Type {
	case "broadcast":
		h.handleBroadcast(node, msg)
	case "read":
		h.handleRead(node, msg)
	case "topology":
		h.handleTopology(node, msg)
	}
}

func (h *BroadcastHandler) handleBroadcast(node *maelstrom.Node, msg *maelstrom.Message) {
	h.messages[*msg.Body.Message] = true
	node.SendMessage(&maelstrom.Message{
		Body: maelstrom.MessageBody{
			Type: "broadcast_ok",
		},
	}, msg, nil)
}

func (h *BroadcastHandler) handleRead(node *maelstrom.Node, msg *maelstrom.Message) {
	messages := maps.Keys(h.messages)
	node.SendMessage(&maelstrom.Message{
		Body: maelstrom.MessageBody{
			Type:     "read_ok",
			Messages: &messages,
		},
	}, msg, nil)
}

func (h *BroadcastHandler) handleTopology(node *maelstrom.Node, msg *maelstrom.Message) {
	h.topology = *msg.Body.Topology
	node.SendMessage(&maelstrom.Message{
		Body: maelstrom.MessageBody{
			Type: "topology_ok",
		},
	}, msg, nil)
}
