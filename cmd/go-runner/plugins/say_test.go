package plugins

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSay(t *testing.T) {
	in := []byte(`{"body":"hello"}`)
	say := &Say{}
	conf, err := say.ParseConf(in)
	assert.Nil(t, err)
	assert.Equal(t, "hello", conf.(SayConf).Body)

	w := httptest.NewRecorder()
	say.Filter(conf, w, nil)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "Go", resp.Header.Get("X-Resp-A6-Runner"))
	assert.Equal(t, "hello", string(body))
}

func TestSay_BadConf(t *testing.T) {
	in := []byte(``)
	say := &Say{}
	_, err := say.ParseConf(in)
	assert.NotNil(t, err)
}

func TestSay_NoBody(t *testing.T) {
	in := []byte(`{}`)
	say := &Say{}
	conf, err := say.ParseConf(in)
	assert.Nil(t, err)
	assert.Equal(t, "", conf.(SayConf).Body)

	w := httptest.NewRecorder()
	say.Filter(conf, w, nil)
	resp := w.Result()
	assert.Equal(t, "", resp.Header.Get("X-Resp-A6-Runner"))
}
