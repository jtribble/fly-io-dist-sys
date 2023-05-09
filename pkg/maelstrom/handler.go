package maelstrom

type Handler interface {
	HandlesMessageType(msgType string) bool

	HandleMessage(node *Node, msg *Message)
}
