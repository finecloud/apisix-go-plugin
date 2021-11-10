package server

import (
	"fmt"

	"github.com/ReneKroon/ttlcache/v2"
	A6Err "github.com/api7/ext-plugin-proto/go/A6/Err"
	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/finecloud/apisix-oauth2-plugin/internal/util"
)

type UnknownType struct {
	ty byte
}

func (err UnknownType) Error() string {
	return fmt.Sprintf("unknown type %d", err.ty)
}

func ReportError(err error) *flatbuffers.Builder {
	builder := util.GetBuilder()
	A6Err.RespStart(builder)

	var code A6Err.Code
	switch err {
	case ttlcache.ErrNotFound:
		code = A6Err.CodeCONF_TOKEN_NOT_FOUND
	default:
		switch err.(type) {
		case UnknownType:
			code = A6Err.CodeBAD_REQUEST
		default:
			code = A6Err.CodeSERVICE_UNAVAILABLE
		}
	}

	A6Err.RespAddCode(builder, code)
	resp := A6Err.RespEnd(builder)
	builder.Finish(resp)
	return builder
}
