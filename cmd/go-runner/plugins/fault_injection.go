package plugins

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"

	pkgHTTP "github.com/finecloud/apisix-oauth2-plugin/pkg/http"
	"github.com/finecloud/apisix-oauth2-plugin/pkg/log"
	"github.com/finecloud/apisix-oauth2-plugin/pkg/plugin"
)

const (
	plugin_name = "fault-injection"
)

func init() {
	err := plugin.RegisterPlugin(&FaultInjection{})
	if err != nil {
		log.Fatalf("failed to register plugin %s: %s", plugin_name, err)
	}
}

// FaultInjection is used in the benchmark
type FaultInjection struct {
}

type FaultInjectionConf struct {
	Body       string `json:"body"`
	HttpStatus int    `json:"http_status"`
	Percentage int    `json:"percentage"`
}

func (p *FaultInjection) Name() string {
	return plugin_name
}

func (p *FaultInjection) ParseConf(in []byte) (interface{}, error) {
	conf := FaultInjectionConf{Percentage: -1}
	err := json.Unmarshal(in, &conf)
	if err != nil {
		return nil, err
	}

	// schema check
	if conf.HttpStatus < 200 {
		return nil, errors.New("bad http_status")
	}
	if conf.Percentage == -1 {
		conf.Percentage = 100
	} else if conf.Percentage < 0 || conf.Percentage > 100 {
		return nil, errors.New("bad percentage")
	}

	return conf, err
}

func sampleHit(percentage int) bool {
	return rand.Intn(100) < percentage
}

func (p *FaultInjection) Filter(conf interface{}, w http.ResponseWriter, r pkgHTTP.Request) {
	fc := conf.(FaultInjectionConf)
	if !sampleHit(fc.Percentage) {
		return
	}

	w.WriteHeader(fc.HttpStatus)
	body := fc.Body
	if len(body) == 0 {
		return
	}

	_, err := w.Write([]byte(body))
	if err != nil {
		log.Errorf("failed to write: %s", err)
	}
}
