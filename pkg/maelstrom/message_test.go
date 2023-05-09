package maelstrom

import (
	"os"
	"testing"

	"github.com/jtribble/fly-io-dist-sys/pkg/ptr"
	"github.com/stretchr/testify/assert"
)

func readFile(t *testing.T, name string) []byte {
	bytes, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}

func Test_From(t *testing.T) {
	msg, err := From(readFile(t, "testdata/init.json"))
	assert.Nil(t, err)
	assert.Equal(t, Message{
		Src:  "c0",
		Dest: "n0",
		Body: MessageBody{
			Type:    "init",
			NodeId:  ptr.ToString("n0"),
			NodeIds: ptr.ToStringSlice([]string{"n0"}),
			MsgId:   ptr.ToInt(1),
		},
	}, *msg)
}
