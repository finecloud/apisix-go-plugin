package util

import (
	"fmt"
	"io"

	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/finecloud/apisix-oauth2-plugin/pkg/log"
)

const (
	HeaderLen   = 4
	MaxDataSize = 2<<24 - 1
)

const (
	RPCError = iota
	RPCPrepareConf
	RPCHTTPReqCall
	RPCExtraInfo
)

type RPCResult struct {
	Err     error
	Builder *flatbuffers.Builder
}

// Use struct if the result is not only []byte
type ExtraInfoResult []byte

func ReadErr(n int, err error, required int) bool {
	if 0 < n && n < required {
		err = fmt.Errorf("truncated, only get the first %d bytes", n)
	}
	if err != nil {
		if err != io.EOF {
			log.Errorf("read: %s", err)
		}
		return true
	}
	return false
}

func WriteErr(n int, err error) {
	if err != nil {
		// TODO: solve "write: broken pipe" with context
		log.Errorf("write: %s", err)
	}
}
