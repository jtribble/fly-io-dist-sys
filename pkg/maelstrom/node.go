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
	id               string
	peers            []string
	handlers         []Handler
	incomingMessages chan *Message
	outgoingMessages chan *Message
	inFlightMessages sync.WaitGroup
	nextMessageId    int
}

func NewNode(handlers ...Handler) *Node {
	handlers = append([]Handler{&InitHandler{}}, handlers...)
	return &Node{
		incomingMessages: make(chan *Message),
		outgoingMessages: make(chan *Message),
		handlers:         handlers,
		nextMessageId:    1,
		inFlightMessages: sync.WaitGroup{},
	}
}

func (n *Node) Id() string {
	return n.id
}

func (n *Node) SendMessage(msg, inReplyTo *Message, _ func(response *Message)) {
	msg.Src = n.id
	if inReplyTo != nil {
		msg.Dest = inReplyTo.Src
		msg.Body.InReplyTo = inReplyTo.Body.MsgId
	}
	msg.Body.MsgId = ptr.ToInt(n.nextMessageId)
	n.nextMessageId += 1
	n.inFlightMessages.Add(1)
	n.outgoingMessages <- msg
}

func (n *Node) RunUntilInterrupted() {
	osutils.RunUntilInterrupted(func(ctx context.Context, cancel context.CancelFunc) {
		go n.handleIncomingMessages(ctx, cancel)
		go n.handleOutgoingMessages()

		<-ctx.Done()

		// Stop accepting incoming messages
		close(n.incomingMessages)

		// Wait a while for outgoing messages to finish
		syncutils.WaitUntilTimeout(&n.inFlightMessages, 10*time.Second)
		close(n.outgoingMessages)
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
			for _, h := range n.handlers {
				if h.HandlesMessageType(msg.Body.Type) {
					h.HandleMessage(n, msg)
					handled = true
				}
			}
			if !handled {
				n.SendMessage(&Message{
					Body: MessageBody{
						Type: "error",
						Code: errors.NotSupported,
					},
				}, msg, nil)
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
		case msg := <-n.outgoingMessages:
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
