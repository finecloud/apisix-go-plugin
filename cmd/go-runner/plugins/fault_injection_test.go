package plugins

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFaultInjection(t *testing.T) {
	in := []byte(`{"http_status":400, "body":"hello"}`)
	fi := &FaultInjection{}
	conf, err := fi.ParseConf(in)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	fi.Filter(conf, w, nil)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "hello", string(body))
}

func TestFaultInjection_Percentage(t *testing.T) {
	in := []byte(`{"http_status":400, "percentage":0}`)
	fi := &FaultInjection{}
	conf, err := fi.ParseConf(in)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	fi.Filter(conf, w, nil)
	resp := w.Result()
	assert.Equal(t, 200, resp.StatusCode)
}
