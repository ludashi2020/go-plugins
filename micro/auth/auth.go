// Package auth provides auth management for micro
package auth

import (
	"net/http"
	"strings"

	"github.com/micro/cli/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/micro/v2/plugin"

	// enterprise auth
	"github.com/ludashi2020/go-plugins/micro/auth/basic"
	"github.com/ludashi2020/go-plugins/micro/auth/digest"
	"github.com/ludashi2020/go-plugins/micro/auth/ldap"
)

type Auth struct {
	Provider Provider
}

// Provider is an auth provider
type Provider interface {
	Handler(h http.Handler) http.Handler
}

// allow always returns success
type allow struct{}

func (a *allow) Handler(h http.Handler) http.Handler {
	return h
}

// Authent
func (a *Auth) Handler(h http.Handler) http.Handler {
	return a.Provider.Handler(h)
}

// NewPlugin returns a new auth plugin
func NewPlugin() plugin.Plugin {
	auth := new(Auth)
	auth.Provider = new(allow)

	return plugin.NewPlugin(
		plugin.WithName("auth"),
		plugin.WithFlag(
			&cli.StringFlag{
				Name:  "auth",
				Usage: "Specify the type of auth e.g basic:///path/to/file, digest:///path/to/file, ldap[s]://url",
			},
			&cli.StringFlag{
				Name:  "realm",
				Usage: "Specify the realm for auth",
			},
		),
		plugin.WithHandler(auth.Handler),
		plugin.WithInit(func(ctx *cli.Context) error {
			authType := ctx.String("auth")
			authRealm := ctx.String("realm")
			parts := strings.Split(authType, "://")

			// no auth
			if len(parts) < 2 {
				return nil
			}

			typ := parts[0]
			file := parts[1]

			switch typ {
			case "basic":
				auth.Provider = basic.New(file, authRealm)
				log.Infof("Loaded basic auth file: %s", file)
			case "digest":
				log.Infof("Loaded digest auth file: %s", file)
				auth.Provider = digest.New(file, authRealm)
			case "ldap", "ldaps":
				log.Infof("Loaded ldap auth url: %s", file)
				auth.Provider = ldap.New(file, authRealm)
			}

			return nil
		}),
	)
}
