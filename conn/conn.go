package conn

import (
	"fmt"

	"github.com/pierelucas/atlantr-extreme/utils"
	"github.com/pierelucas/gorpc"
)

// Send the payload to the backend server
func Send(payload string, backend string) error {
	encodedPayload := utils.Base64Encode(payload)

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

	if resp.(int) != 0 {
		err := fmt.Errorf("There was an error on Serverside: %+v", resp)
		return err
	}

	return nil
}
