package uploader

import (
	"encoding/json"

	"github.com/pierelucas/atlantr-extreme/data"
)

// Pair --
type Pair struct {
	ID            data.Value
	EMAILUSER     data.Value
	EMAILPASSWORD data.Value
}

// NewPair generates a new pair and calls utils.GenerateID
func NewPair(user, pass string, id string) (*Pair, error) {
	return &Pair{
		ID:            data.Value(id),
		EMAILUSER:     data.Value(user),
		EMAILPASSWORD: data.Value(pass),
	}, nil
}

// GetUser gets the email user
func (p *Pair) GetUser() string {
	return p.EMAILUSER.String()
}

// GetPassword gets the email password
func (p *Pair) GetPassword() string {
	return p.EMAILPASSWORD.String()
}

// SetUser sets the email user
func (p *Pair) SetUser(user string) {
	p.EMAILUSER = data.Value(user)
}

// SetPassword sets the email password
func (p *Pair) SetPassword(pass string) {
	p.EMAILPASSWORD = data.Value(pass)
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
