package uploader

import (
	"encoding/json"

	"github.com/pierelucas/atlantr-extreme/data"
)

// Pair --
type Pair struct {
	ID           data.Value
	VALIDDATA    data.Value
	NOTFOUNDDATA data.Value
}

// NewPair generates a new pair and calls utils.GenerateID
func NewPair(vdata, nfdata []byte, id string) (*Pair, error) {
	return &Pair{
		ID:           data.Value(id),
		VALIDDATA:    data.Value(vdata),
		NOTFOUNDDATA: data.Value(nfdata),
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
