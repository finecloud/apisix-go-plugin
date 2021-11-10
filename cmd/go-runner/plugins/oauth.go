package plugins

import (
	"encoding/json"
	"fmt"
	pkgHTTP "github.com/finecloud/apisix-oauth2-plugin/pkg/http"
	"github.com/finecloud/apisix-oauth2-plugin/pkg/log"
	"github.com/finecloud/apisix-oauth2-plugin/pkg/plugin"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func init() {
	err := plugin.RegisterPlugin(&Oauth{})
	if err != nil {
		log.Fatalf("failed to register plugin OauthConf: %s", err)
	}
}

// Oauth is for spring oauth
type Oauth struct {
}

type OauthConf struct {
	ApiKey   string `json:"api_key"`
	Password string `json:"password"`
	CheckUrl string `json:"check_url"`
}

func (p *Oauth) Name() string {
	return "Oauth2"
}

func (p *Oauth) ParseConf(in []byte) (interface{}, error) {
	conf := OauthConf{}
	err := json.Unmarshal(in, &conf)
	return conf, err
}

func (p *Oauth) Filter(conf interface{}, w http.ResponseWriter, r pkgHTTP.Request) {
	oauthConf := conf.(OauthConf)
	apiKey := oauthConf.ApiKey
	password := oauthConf.Password
	checkUrl := oauthConf.CheckUrl

	if len(apiKey) == 0 {
		log.Infof("oauth plugin not find conf for api_key")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if len(password) == 0 {
		log.Infof("oauth plugin not find conf for password")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if len(checkUrl) == 0 {
		log.Infof("oauth plugin not find conf for check_url")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	authorization := r.Header().Get("Authorization")
	if len(authorization) == 0 {
		log.Infof("request header not find authorization")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !strings.Contains(authorization, "Bearer") {
		log.Infof("error authorization %s", authorization)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := strings.Replace(authorization, "Bearer ", "", 1)
	isToken := verifyToken(token, apiKey, password, checkUrl)
	if isToken {
		return
	} else {
		log.Infof("Illegal token %s", authorization)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

type oauthCheckResp struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// verifyToken Verify that the token is correct
func verifyToken(token, apiKey, password, checkUrl string) bool {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?token=%s", checkUrl, token), nil)
	if err != nil {
		panic(err.Error())
	}
	req.SetBasicAuth(apiKey, password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	checkResp := oauthCheckResp{}
	err = json.Unmarshal(body, &checkResp)
	if err != nil {
		return false
	}
	if len(checkResp.Error) != 0 || checkResp.Error == "invalid_token" {
		return false
	}
	return true
}
