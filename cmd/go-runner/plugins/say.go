package plugins

import (
	"encoding/json"
	"net/http"

	pkgHTTP "github.com/finecloud/apisix-oauth2-plugin/pkg/http"
	"github.com/finecloud/apisix-oauth2-plugin/pkg/log"
	"github.com/finecloud/apisix-oauth2-plugin/pkg/plugin"
)

func init() {
	err := plugin.RegisterPlugin(&Say{})
	if err != nil {
		log.Fatalf("failed to register plugin say: %s", err)
	}
}

// Say is a demo to show how to return data directly instead of proxying
// it to the upstream.
type Say struct {
}

type SayConf struct {
	Body string `json:"body"`
}

func (p *Say) Name() string {
	return "say"
}

func (p *Say) ParseConf(in []byte) (interface{}, error) {
	conf := SayConf{}
	err := json.Unmarshal(in, &conf)
	return conf, err
}

func (p *Say) Filter(conf interface{}, w http.ResponseWriter, r pkgHTTP.Request) {
	body := conf.(SayConf).Body
	if len(body) == 0 {
		return
	}

	w.Header().Add("X-Resp-A6-Runner", "Go")
	_, err := w.Write([]byte(body))
	if err != nil {
		log.Errorf("failed to write: %s", err)
	}
}
