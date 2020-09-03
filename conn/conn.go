package conn

import (
	"encoding/base64"
	"fmt"

	"github.com/pierelucas/gorpc"
)

func base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Send(payload string, backend string) error {
	encodedPayload := base64Encode(payload)

	// RPC Client
	c := &gorpc.Client{
		// TCP address of the server.
		Addr: backend,
	}

	// start client
	c.Start()

	// send payload
	resp, err := c.Call(encodedPayload)
	if err != nil {
		return err
	}
	if resp.(string) != encodedPayload {
		err := fmt.Errorf("Unexpected response from the server: %+v", resp)
		return err
	}

	return nil
}
