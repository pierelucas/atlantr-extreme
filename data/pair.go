package data

import (
	"encoding/json"

	"github.com/pierelucas/atlantr-extreme/conn"
	"github.com/pierelucas/atlantr-extreme/utils"
)

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
	jsonClient, err := func() (string, error) {
		data, err := json.Marshal(p)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}()
	if err != nil {
		return "", err
	}

	return jsonClient, nil
}

// SendToServer Marshals the structure to JSON and calls conn.Send(). The Payload is then b64 encoded and send to the backend.
// In case that an error occur, this error will be returned
func (p *Pair) SendToServer(backend string) error {
	jsonPair, err := func() (string, error) {
		data, err := json.Marshal(p)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}()
	if err != nil {
		return err
	}

	return conn.Send(jsonPair, backend)
}
