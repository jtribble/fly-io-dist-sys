package maelstrom

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/jtribble/fly-io-dist-sys/pkg/ioutils"
	"github.com/jtribble/fly-io-dist-sys/pkg/log"
	"github.com/jtribble/fly-io-dist-sys/pkg/maelstrom/errors"
	"github.com/jtribble/fly-io-dist-sys/pkg/osutils"
	"github.com/jtribble/fly-io-dist-sys/pkg/ptr"
	"github.com/jtribble/fly-io-dist-sys/pkg/syncutils"
)

type Node struct {
	Id    string
	Peers []string

	Initialized bool

	IncomingMessages chan *Message
	OutgoingMessages chan *Message

	Handlers []Handler

	nextMessageId    int
	inFlightMessages sync.WaitGroup
}

func NewNode() *Node {
	return &Node{
		IncomingMessages: make(chan *Message),
		OutgoingMessages: make(chan *Message),
		Handlers:         []Handler{&InitHandler{}},
		nextMessageId:    1,
		inFlightMessages: sync.WaitGroup{},
	}
}

func (n *Node) RegisterHandler(h Handler) {
	n.Handlers = append(n.Handlers, h)
}

func (n *Node) QueueReply(msg, inReplyTo *Message) {
	msg.Src = n.Id
	msg.Dest = inReplyTo.Src
	msg.Body.InReplyTo = inReplyTo.Body.MsgId
	msg.Body.MsgId = ptr.ToInt(n.nextMessageId)
	n.nextMessageId += 1
	n.inFlightMessages.Add(1)
	n.OutgoingMessages <- msg
}

func (n *Node) RunUntilInterrupted() {
	osutils.RunUntilInterrupted(func(ctx context.Context, cancel context.CancelFunc) {
		go n.handleIncomingMessages(ctx, cancel)
		go n.handleOutgoingMessages()

		<-ctx.Done()

		// Stop accepting incoming messages
		close(n.IncomingMessages)

		// Wait a while for outgoing messages to finish
		syncutils.WaitUntilTimeout(&n.inFlightMessages, 10*time.Second)
		close(n.OutgoingMessages)
	})
}

func (n *Node) handleIncomingMessages(ctx context.Context, cancel context.CancelFunc) {
	lr := ioutils.NewLineReader(bufio.NewReader(os.Stdin))
	for {
		select {
		case <-ctx.Done():
			return
		case line := <-lr.Lines:
			msg, err := From(line)
			if err != nil {
				log.Stderrf("failed to parse message: %s", err)
				continue
			}
			handled := false
			for _, h := range n.Handlers {
				if h.HandlesMessageType(msg.Body.Type) {
					h.HandleMessage(n, msg)
					handled = true
				}
			}
			if !handled {
				n.QueueReply(&Message{
					Body: MessageBody{
						Type: "error",
						Code: errors.NotSupported,
					},
				}, msg)
			}
		case err := <-lr.Errs:
			log.Stderrf("failed to read from stdin: %s", err)
		case <-lr.Eof:
			cancel()
			return
		}
	}
}

func (n *Node) handleOutgoingMessages() {
	for {
		select {
		case msg := <-n.OutgoingMessages:
			n.writeMessageToStdout(msg)
		}
	}
}

func (n *Node) writeMessageToStdout(msg *Message) {
	defer n.inFlightMessages.Done()
	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Stderrf("failed to serialize message: %s", err)
		return
	}
	log.Stdout(string(bytes))
}
