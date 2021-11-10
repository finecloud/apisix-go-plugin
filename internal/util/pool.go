package util

import (
	"sync"

	flatbuffers "github.com/google/flatbuffers/go"
)

var builderPool = sync.Pool{
	New: func() interface{} {
		return flatbuffers.NewBuilder(256)
	},
}

func GetBuilder() *flatbuffers.Builder {
	return builderPool.Get().(*flatbuffers.Builder)
}

func PutBuilder(b *flatbuffers.Builder) {
	b.Reset()
	builderPool.Put(b)
}
