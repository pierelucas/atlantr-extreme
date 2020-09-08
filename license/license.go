package license

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
