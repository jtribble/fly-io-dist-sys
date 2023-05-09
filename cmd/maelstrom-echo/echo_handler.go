package main

import "github.com/jtribble/fly-io-dist-sys/pkg/maelstrom"

type EchoHandler struct {
}

func (h EchoHandler) HandlesMessageType(msgType string) bool {
	return msgType == "echo"
}

func (h EchoHandler) HandleMessage(node *maelstrom.Node, msg *maelstrom.Message) {
	node.QueueReply(&maelstrom.Message{
		Body: maelstrom.MessageBody{
			Type: "echo_ok",
			Echo: msg.Body.Echo,
		},
	}, msg)
}
