package maelstrom

type InitHandler struct {
}

func (h InitHandler) HandlesMessageType(msgType string) bool {
	return msgType == "init"
}

func (h InitHandler) HandleMessage(node *Node, msg *Message) {
	if !node.Initialized {
		node.Id = *msg.Body.NodeId
		node.Peers = *msg.Body.NodeIds
		node.QueueReply(&Message{
			Body: MessageBody{
				Type: "init_ok",
			},
		}, msg)
	}
	node.Initialized = true
}
