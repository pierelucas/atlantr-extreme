package data

import (
	"fmt"
)

// Host struct holding data for host addr and port
type Host struct {
	host Value
	port Value
}

// NewHost returns Host pointer
func NewHost(host, port string) *Host {
	return &Host{
		host: Value(host),
		port: Value(port),
	}
}

// GetHost --
func (h *Host) GetHost() string {
	return h.host.String()
}

// GetPort --
func (h *Host) GetPort() string {
	return h.port.String()
}

// GetFullAddr of Host:Port as string
func (h *Host) GetFullAddr() string {
	return fmt.Sprintf("%s:%s", h.GetHost(), h.GetPort())
}
