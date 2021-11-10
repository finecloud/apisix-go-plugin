package plugin

import (
	"net/http"

	"github.com/finecloud/apisix-oauth2-plugin/internal/plugin"
	pkgHTTP "github.com/finecloud/apisix-oauth2-plugin/pkg/http"
)

// Plugin represents the Plugin
type Plugin interface {
	// Name returns the plguin name
	Name() string

	// ParseConf is the method to parse given plugin configuration. When the
	// configuration can't be parsed, it will be skipped.
	ParseConf(in []byte) (conf interface{}, err error)

	// Filter is the method to handle request.
	// It is like the `http.ServeHTTP`, plus the ctx and the configuration created by
	// ParseConf.
	//
	// When the `w` is written, the execution of plugin chain will be stopped.
	// We don't use onion model like Gin/Caddy because we don't serve the whole request lifecycle
	// inside the runner. The plugin is only a filter running at one stage.
	Filter(conf interface{}, w http.ResponseWriter, r pkgHTTP.Request)
}

// RegisterPlugin register a plugin. Plugin which has the same name can't be registered twice.
// This method should be called before calling `runner.Run`.
func RegisterPlugin(p Plugin) error {
	return plugin.RegisterPlugin(p.Name(), p.ParseConf, p.Filter)
}
