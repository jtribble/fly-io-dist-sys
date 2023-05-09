package log

import (
	"fmt"
	"os"
)

func Stdout(a ...any) {
	_, _ = fmt.Fprintln(os.Stdout, a...)
}

func Stdoutf(format string, a ...any) {
	Stderr(fmt.Sprintf(format, a...))
}

func Stderr(a ...any) {
	_, _ = fmt.Fprintln(os.Stderr, a...)
}

func Stderrf(format string, a ...any) {
	Stderr(fmt.Sprintf(format, a...))
}
