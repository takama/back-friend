package handlers

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/takama/back-friend/pkg/version"

	"github.com/takama/bit"
)

// Status contains detailed information about service
type Status struct {
	Host     string   `json:"host"`
	Version  string   `json:"version"`
	Commit   string   `json:"commit"`
	Repo     string   `json:"repo"`
	Compiler string   `json:"compiler"`
	Runtime  Runtime  `json:"runtime"`
	State    State    `json:"state"`
	Requests Requests `json:"requests"`
}

// Runtime defines runtime part of service information
type Runtime struct {
	CPU        int    `json:"cpu"`
	Memory     string `json:"memory"`
	Goroutines int    `json:"goroutines"`
}

// State contains current state of the service
type State struct {
	Maintenance bool   `json:"maintenance"`
	Uptime      string `json:"uptime"`
}

// Requests contains response codes statistics
type Requests struct {
	C2xx int `json:"2xx"`
	C4xx int `json:"4xx"`
	C5xx int `json:"5xx"`
}

// Info returns detailed info about the service
func (h *Handler) Info(c bit.Control) {
	host, _ := os.Hostname()
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)

	c.Code(http.StatusOK)
	c.Body(Status{
		Host:     host,
		Version:  version.RELEASE,
		Commit:   version.COMMIT,
		Repo:     version.REPO,
		Compiler: runtime.Version(),
		Runtime: Runtime{
			CPU:        runtime.NumCPU(),
			Memory:     fmt.Sprintf("%.2fMB", float64(m.Sys)/(1<<(10*2))),
			Goroutines: runtime.NumGoroutine(),
		},
		State: State{
			Maintenance: h.maintenance,
			Uptime:      time.Now().Sub(h.stats.startTime).String(),
		},
		Requests: Requests{
			C2xx: h.stats.requests.C2xx,
			C4xx: h.stats.requests.C4xx,
			C5xx: h.stats.requests.C5xx,
		},
	})
}
