package socket

import (
	"bytes"
	"context"
	"testing"

	waLog "github.com/polymorfa/hypermeow/util/log"
)

func TestFrameStreamBoundaries(t *testing.T) {
	stream := []byte{0, 0, 2, 1, 2, 0, 0, 0, 0, 0, 3, 3, 4, 5}
	want := [][]byte{{1, 2}, {}, {3, 4, 5}}
	// Split at every pair of boundaries, including inside a header or payload.
	for i := 0; i <= len(stream); i++ {
		for j := i; j <= len(stream); j++ {
			fs := NewFrameSocket(waLog.Noop, nil)
			fs.cancelCtx = context.Background()
			fs.Frames = make(chan []byte, len(want)+1)
			fs.processData(stream[:i:i])
			fs.processData(stream[i:j:j])
			fs.processData(stream[j:])
			if len(fs.Frames) != len(want) {
				t.Fatalf("splits %d/%d yielded %d frames", i, j, len(fs.Frames))
			}
			for _, v := range want {
				if got := <-fs.Frames; !bytes.Equal(got, v) {
					t.Fatalf("splits %d/%d yielded %v", i, j, got)
				}
			}
		}
	}
}

func BenchmarkSplitFrameHeader(b *testing.B) {
	fs := NewFrameSocket(waLog.Noop, nil)
	fs.cancelCtx = context.Background()
	fs.Frames = make(chan []byte, 1)
	head := []byte{1}
	rest := make([]byte, 2+65536)
	b.ReportAllocs()
	for b.Loop() {
		fs.processData(head)
		fs.processData(rest)
		<-fs.Frames
	}
}
