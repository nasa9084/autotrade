package bufferpool_test

import (
	"fmt"
	"testing"

	"github.com/nasa9084/autotrade/bufferpool"
)

func TestBufferPool(t *testing.T) {
	const s = `blahblahblah`

	buf := bufferpool.Get()

	buf.WriteString(s)
	if got := buf.String(); got != s {
		t.Fatalf("unexpected buffer content: %s != %s", got, s)
	}

	if got := buf.Len(); got != 12 { // len(blahblahblah) = 12
		t.Fatalf("unexpected buffer length: %d != 12", got)
	}

	bufferpool.Put(buf)

	if got := buf.Len(); got != 0 {
		t.Fatalf("unexpected buffer length: %d != 0", got)
	}

	buf2 := bufferpool.Get()

	p1 := fmt.Sprintf("%p", buf)
	p2 := fmt.Sprintf("%p", buf2)

	if p1 != p2 {
		t.Fatalf("unexpected buffer address: %s != %s", p2, p1)
	}

	buf3 := bufferpool.Get() // should be new one
	p3 := fmt.Sprintf("%p", buf3)

	if p1 == p3 {
		t.Fatalf("unexpected buffer address: %s != %s", p3, p1)
	}
}
