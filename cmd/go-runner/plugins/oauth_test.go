package plugins

import (
	pkgHTTP "github.com/finecloud/apisix-oauth2-plugin/pkg/http"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type TestHead struct {
	Data map[string]string
}

func (r *TestHead) Set(key, value string) {
	r.Data[key] = value
}

func (r *TestHead) Del(key string) {

}

func (r *TestHead) Get(key string) string {
	return r.Data[key]
}

func (r *TestHead) View() http.Header {
	return nil
}

type TestRequest struct {
	Head *TestHead
}

func (r *TestRequest) ID() uint32 {
	return 1
}
func (r *TestRequest) SrcIP() net.IP {
	return nil
}
func (r *TestRequest) Method() string {
	return ""

}
func (r *TestRequest) Path() []byte {
	return nil

}
func (r *TestRequest) SetPath([]byte) {
	return
}
func (r *TestRequest) Header() pkgHTTP.Header {
	if r.Head == nil {
		r.Head = &TestHead{
			Data: map[string]string{},
		}
	}
	return r.Head
}
func (r *TestRequest) Args() url.Values {
	return nil

}
func (r *TestRequest) Var(name string) ([]byte, error) {
	return nil, nil

}

func TestOauth2Fail(t *testing.T) {
	in := []byte(`{"api_key":"app","password":"app","check_url":"http://dev.com:9999/auth/oauth/check_token"}`)
	oauth := &Oauth{}
	conf, err := oauth.ParseConf(in)
	assert.Nil(t, err)
	assert.Equal(t, "app", conf.(OauthConf).ApiKey)

	w := httptest.NewRecorder()
	request := &TestRequest{}
	request.Header().Set("Authorization", "Bearer 6a5a8f05-34a9-4a83-8deb-012892c8b08f")
	oauth.Filter(conf, w, request)
	resp := w.Result()
	assert.Equal(t, 401, resp.StatusCode)
}

func TestOauth2Success(t *testing.T) {
	in := []byte(`{"api_key":"app","password":"app","check_url":"http://dev.com:9999/auth/oauth/check_token"}`)
	oauth := &Oauth{}
	conf, err := oauth.ParseConf(in)
	assert.Nil(t, err)
	assert.Equal(t, "app", conf.(OauthConf).ApiKey)

	w := httptest.NewRecorder()
	request := &TestRequest{}
	request.Header().Set("Authorization", "Bearer 9fa7c1ac-4a29-48b2-83e8-d03be43f0ac6")
	oauth.Filter(conf, w, request)
	resp := w.Result()
	//body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 200, resp.StatusCode)
}
