package plugins

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	pkgHTTP "github.com/finecloud/apisix-oauth2-plugin/pkg/http"
	"github.com/finecloud/apisix-oauth2-plugin/pkg/log"
	"github.com/finecloud/apisix-oauth2-plugin/pkg/plugin"
)

func init() {
	err := plugin.RegisterPlugin(&LimitReq{})
	if err != nil {
		log.Fatalf("failed to register plugin limit-req: %s", err)
	}
}

// LimitReq is a demo for a real world plugin
type LimitReq struct {
}

type LimitReqConf struct {
	Burst int     `json:"burst"`
	Rate  float64 `json:"rate"`

	limiter *rate.Limiter
}

func (p *LimitReq) Name() string {
	return "limit-req"
}

// ParseConf is called when the configuration is changed. And its output is unique per route.
func (p *LimitReq) ParseConf(in []byte) (interface{}, error) {
	conf := LimitReqConf{}
	err := json.Unmarshal(in, &conf)
	if err != nil {
		return nil, err
	}

	limiter := rate.NewLimiter(rate.Limit(conf.Rate), conf.Burst)
	// the conf can be used to store route scope data
	conf.limiter = limiter
	return conf, nil
}

// Filter is called when a request hits the route
func (p *LimitReq) Filter(conf interface{}, w http.ResponseWriter, r pkgHTTP.Request) {
	li := conf.(LimitReqConf).limiter
	rs := li.Reserve()
	if !rs.OK() {
		// limit rate exceeded
		log.Infof("limit req rate exceeded")
		// stop filters with this response
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	time.Sleep(rs.Delay())
}
