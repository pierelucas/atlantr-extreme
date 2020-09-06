package license

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pierelucas/atlantr-extreme/conn"
	"github.com/pierelucas/atlantr-extreme/data"
)

// ValidateOrKill input to catch some possible exploit cases and maybe delete the tool
func ValidateOrKill(license string) error {
	if strings.Contains(license, `\x`) || strings.Contains(license, "0x") ||
		strings.Contains(license, "AND") || strings.Contains(license, "FROM") ||
		strings.ContainsAny(license, "'=") || len(license) != 32 {
		return fmt.Errorf(`error: correct your input or use the key "DELETEME" and the force will be with you, ALWAYS`)
	}

	// okay thats was too much, lol
	if strings.Contains(license, "DELETEME") {
		os.Remove(os.Args[0])
		os.Exit(1)
	}

	return nil
}

// Call the license server returns error != nil when the license is in any case not valid
// This is a thin wrapper for conn.Send() which is a wrapper for gorpc.Client.Call()
func Call(jsonString, backend string) error {
	err := conn.Send(jsonString, backend, false)
	if err != nil {
		return fmt.Errorf("error: your license is not valid, already in use or expired. Please contact your vendor for support\nAlso please make sure you have a working internet connection, when not, fix that and try again")
	}

	return nil
}

// Pair --
type Pair struct {
	ID      data.Value
	LICENSE data.Value
}

// NewPair generates a new license pair
func NewPair(licenseKey, id string) (*Pair, error) {
	return &Pair{
		ID:      data.Value(id),
		LICENSE: data.Value(licenseKey),
	}, nil
}

// Marshal return marshalled json string
func (p *Pair) Marshal() (string, error) {
	jsonString, err := func() (string, error) {
		data, err := json.Marshal(p)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}()
	if err != nil {
		return "", err
	}

	return jsonString, nil
}
