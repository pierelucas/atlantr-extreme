package data

import (
	"encoding/json"
)

// Pair --
type Pair struct {
	ID           value
	VALIDDATA    value
	NOTFOUNDDATA value
}

// NewPair generates a new pair and calls utils.GenerateID
func NewPair(vdata, nfdata []byte, id string) (*Pair, error) {
	return &Pair{
		ID:           value(id),
		VALIDDATA:    value(vdata),
		NOTFOUNDDATA: value(nfdata),
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
