package bufferpool

import (
	"bytes"
	"sync"
)

var global = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func Get() *bytes.Buffer {
	return global.Get().(*bytes.Buffer)
}

func Put(buf *bytes.Buffer) {
	buf.Reset()
	global.Put(buf)
}
