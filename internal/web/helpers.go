package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"runtime/debug"
	"strings"
)

func (a *app) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	a.errorLog.Output(2, trace)

	if a.appConfig.Debug {
		http.Error(w, trace, http.StatusInternalServerError)

		return
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *app) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (a *app) forbidden(w http.ResponseWriter) {
	a.clientError(w, http.StatusForbidden)
}

func (a *app) notFound(w http.ResponseWriter) {
	a.clientError(w, http.StatusNotFound)
}

func (a *app) methodNotAllowed(w http.ResponseWriter, allow []string) {
	w.Header().Set("Allow", strings.Join(allow, ", "))
	a.clientError(w, http.StatusMethodNotAllowed)
}

func withoutURL(err error) error {
	if urlErr, ok := errors.AsType[*url.Error](err); ok {
		return urlErr.Err
	}

	return err
}

func requestHostname(host string) string {
	if hostname, _, err := net.SplitHostPort(host); err == nil {
		host = hostname
	}

	return strings.ToLower(strings.TrimSuffix(strings.TrimPrefix(host, "["), "]"))
}

func isLoopback(hostname string) bool {
	if hostname == "localhost" {
		return true
	}

	addr, err := netip.ParseAddr(hostname)

	return err == nil && addr.IsLoopback()
}

func jsonEscape(value string) string {
	encoded, _ := json.Marshal(value)

	return string(encoded[1 : len(encoded)-1])
}
