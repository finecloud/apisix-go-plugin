package server

import (
	"io"
	"testing"
	"time"

	A6Err "github.com/api7/ext-plugin-proto/go/A6/Err"
	"github.com/finecloud/apisix-oauth2-plugin/internal/plugin"
	"github.com/stretchr/testify/assert"
)

func TestReportErrorCacheToken(t *testing.T) {
	plugin.InitConfCache(10 * time.Millisecond)

	_, err := plugin.GetRuleConf(uint32(999999))
	b := ReportError(err)
	out := b.FinishedBytes()
	resp := A6Err.GetRootAsResp(out, 0)
	assert.Equal(t, A6Err.CodeCONF_TOKEN_NOT_FOUND, resp.Code())
}

func TestReportErrorUnknownType(t *testing.T) {
	b := ReportError(UnknownType{23})
	out := b.FinishedBytes()
	resp := A6Err.GetRootAsResp(out, 0)
	assert.Equal(t, A6Err.CodeBAD_REQUEST, resp.Code())
}

func TestReportErrorUnknownErr(t *testing.T) {
	b := ReportError(io.EOF)
	out := b.FinishedBytes()
	resp := A6Err.GetRootAsResp(out, 0)
	assert.Equal(t, A6Err.CodeSERVICE_UNAVAILABLE, resp.Code())
}
