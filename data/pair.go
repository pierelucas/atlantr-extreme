package data

import (
	"encoding/json"

	"github.com/pierelucas/atlantr-extreme/utils"
)

// Pair --
type Pair struct {
	ID           value
	VALIDDATA    value
	NOTFOUNDDATA value
}

// NewPair generates a new pair and calls utils.GenerateID
func NewPair(vdata, nfdata []byte) *Pair {
	// Generate computer ID
	uuid := utils.GenerateID()

	return &Pair{
		ID:           value(uuid),
		VALIDDATA:    value(vdata),
		NOTFOUNDDATA: value(nfdata),
	}
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
