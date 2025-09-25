package log

import (
	"fmt"
	"sync"
)

type Sink interface {
	write(msg message)
}

type stdoutSink struct {
	mu sync.Mutex
}

func (s *stdoutSink) write(msg message) {
	s.mu.Lock()
	// TODO: use os.Stdout
	// n, err = w.Write(p.buf)
	fmt.Println(msg.msg)
	s.mu.Unlock()
}

func StdoutSink() Sink {
	return &stdoutSink{}
}
