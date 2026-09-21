package metrics

import (
	"net/http"
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var ready atomic.Bool

func SetReady(v bool) { ready.Store(v) }
func Ready() bool     { return ready.Load() }

func Handler() http.Handler {
	return promhttp.Handler()
}
