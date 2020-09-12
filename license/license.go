package license

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"github.com/pierelucas/atlantr-extreme/data"
)

// ValidateOrKill input to catch some possible exploit cases and maybe delete the tool
func ValidateOrKill(license string) (string, error) {
	// remove whitespaces, control characters or punctuation signs (e.g minus sign -) from the key
	license = strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) || unicode.IsSpace(r) || unicode.IsControl(r) {
			return -1
		}

		return r
	}, license)

	// check for bad string subsets or the wrong len(license)
	if strings.Contains(license, `\x`) || strings.Contains(license, "0x") ||
		strings.Contains(license, "AND") || strings.Contains(license, "FROM") ||
		strings.ContainsAny(license, "'=") || len(license) != 32 {
		return "", fmt.Errorf("error: license key is not valid! Please delete license.dat and insert valid license key")
	}

	return license, nil
}

// Pair --
type Pair struct {
	ID      data.Value
	LICENSE data.Value
	APPID   data.Value
}

// NewPair generates a new license pair
func NewPair(licenseKey, id, appID string) (*Pair, error) {
	return &Pair{
		ID:      data.Value(id),
		LICENSE: data.Value(licenseKey),
		APPID:   data.Value(appID),
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
