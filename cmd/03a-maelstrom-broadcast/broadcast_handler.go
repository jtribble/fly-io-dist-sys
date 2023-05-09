package main

import (
	"github.com/jtribble/fly-io-dist-sys/pkg/maelstrom"
)

type BroadcastHandler struct {
	messages map[int]bool
	topology *maelstrom.Topology
}

func NewBroadcastHandler() *BroadcastHandler {
	return &BroadcastHandler{
		messages: make(map[int]bool),
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
	node.QueueReply(&maelstrom.Message{
		Body: maelstrom.MessageBody{
			Type: "broadcast_ok",
		},
	}, msg)
}

func (h *BroadcastHandler) handleRead(node *maelstrom.Node, msg *maelstrom.Message) {
	messages := make([]int, len(h.messages))
	i := 0
	for m := range h.messages {
		messages[i] = m
		i += 1
	}
	node.QueueReply(&maelstrom.Message{
		Body: maelstrom.MessageBody{
			Type:     "read_ok",
			Messages: &messages,
		},
	}, msg)
}

func (h *BroadcastHandler) handleTopology(node *maelstrom.Node, msg *maelstrom.Message) {
	h.topology = msg.Body.Topology
	node.QueueReply(&maelstrom.Message{
		Body: maelstrom.MessageBody{
			Type: "topology_ok",
		},
	}, msg)
}
