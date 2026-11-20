package web

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/petaki/inertia-go"
	"github.com/petaki/satellite/internal/config"
	satellitemcp "github.com/petaki/satellite/internal/mcp"
	"github.com/petaki/support-go/cli"
)

func newHostTestApp(allowedHosts ...string) *app {
	return &app{
		appConfig: &config.Config{AllowedHosts: allowedHosts},
		errorLog:  log.New(io.Discard, "", 0),
		infoLog:   log.New(io.Discard, "", 0),
	}
}

func serveHost(handler http.Handler, host string) (int, bool) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = host
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	return rr.Code, rr.Code == http.StatusTeapot
}

func TestAllowedHosts(t *testing.T) {
	cases := map[string]struct {
		allowed []string
		host    string
		want    bool
	}{
		"app url host":             {[]string{"satellite.example"}, "satellite.example", true},
		"app url host on any port": {[]string{"satellite.example"}, "satellite.example:8080", true},
		"case is ignored":          {[]string{"satellite.example"}, "Satellite.EXAMPLE:4000", true},
		"second listed host":       {[]string{"satellite.example", "10.0.0.5"}, "10.0.0.5:4000", true},
		"localhost":                {[]string{"satellite.example"}, "localhost:4000", true},
		"loopback ipv4":            {[]string{"satellite.example"}, "127.0.0.1:4000", true},
		"other loopback ipv4":      {[]string{"satellite.example"}, "127.1.2.3", true},
		"loopback ipv6":            {[]string{"satellite.example"}, "[::1]:4000", true},
		"loopback ipv6 no port":    {[]string{"satellite.example"}, "[::1]", true},
		"wildcard":                 {[]string{"*"}, "attacker.example", true},
		"listed ipv6":              {[]string{"fd00::5"}, "[FD00::5]:4000", true},
		"foreign host":             {[]string{"satellite.example"}, "attacker.example", false},
		"foreign host with port":   {[]string{"satellite.example"}, "attacker.example:4000", false},
		"suffix of allowed host":   {[]string{"satellite.example"}, "evil-satellite.example", false},
		"subdomain of allowed":     {[]string{"satellite.example"}, "x.satellite.example", false},
		"localhost subdomain":      {[]string{"satellite.example"}, "localhost.attacker.example", false},
		"unlisted lan ip":          {[]string{"satellite.example"}, "192.168.0.10:4000", false},
		"empty host":               {[]string{"satellite.example"}, "", false},
		"empty list":               {nil, "attacker.example", false},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			a := newHostTestApp(c.allowed...)
			handler := a.allowedHosts(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusTeapot)
			}))

			code, called := serveHost(handler, c.host)

			if c.want && !called {
				t.Errorf("Host %q with allowed %v: status %d, want the handler to run", c.host, c.allowed, code)
			}

			if !c.want && (called || code != http.StatusForbidden) {
				t.Errorf("Host %q with allowed %v: status %d, handler called %v, want 403 and no call", c.host, c.allowed, code, called)
			}
		})
	}
}

// A rebinding page reaches satellite with its own domain in Host. The guard
// must sit in front of both the destructive web routes and /mcp, even when the
// connection's local address is not loopback (as in Docker).
func TestRoutesRejectForeignHost(t *testing.T) {
	a := newHostTestApp("satellite.example")
	a.sessionManager = scs.New()
	a.inertiaManager = inertia.New("http://satellite.example", "app.gohtml", "v")
	a.mcpServer = satellitemcp.NewHandler(&cli.App{Name: "Satellite", Version: "test"}, a.appConfig, nil, nil, nil, nil)

	routes := a.routes()
	local, _ := net.ResolveTCPAddr("tcp", "172.17.0.2:4000")
	body := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"x","version":"1"}}}`

	for _, target := range []string{"/probe/delete-all", "/mcp"} {
		method := http.MethodDelete

		if target == "/mcp" {
			method = http.MethodPost
		}

		req := httptest.NewRequest(method, target, strings.NewReader(body))
		req.Host = "attacker.example"
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json, text/event-stream")
		req = req.WithContext(context.WithValue(req.Context(), http.LocalAddrContextKey, local))
		rr := httptest.NewRecorder()

		routes.ServeHTTP(rr, req)

		if rr.Code != http.StatusForbidden {
			t.Errorf("%s %s with Host attacker.example: status %d, want 403", method, target, rr.Code)
		}
	}
}
