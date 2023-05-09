package osutils

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jtribble/fly-io-dist-sys/pkg/log"
)

func RunUntilInterrupted(f func(ctx context.Context, cancel context.CancelFunc)) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		log.Stderrf("received signal: %s", <-sigs)
		cancel()
	}()
	f(ctx, cancel)
}
