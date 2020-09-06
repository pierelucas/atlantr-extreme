package conn

import (
	"fmt"
	"log"

	"github.com/pierelucas/atlantr-extreme/utils"
	"github.com/pierelucas/gorpc"
)

// Send the payload to the backend server. This is a one-time rpc call.
// The Client will be terminated automatically from the Go runtime after the return.
// Have in mind that conn.Send() will send the request as base64 encoded string. So the server has to decode the request.
// conn.Send() awaits a response from int = 0 if the call success and int = 1 if the call fails. If the response is int = 1 error is not-nil.
func Send(payload, backend string, logError bool) error {
	encodedPayload := utils.Base64Encode(payload)

	// RPC Client
	c := &gorpc.Client{
		// TCP address of the server.
		Addr: backend,

		//Loggerfunc, we don't use the default func cause the remote address, or other sensible data, is shown in log when a error occurs.
		// So we use logError true if we want to print out the error. In most cases we dont want to print out the error.
		LogError: func(format string, args ...interface{}) {
			if logError {
				log.Printf(format, args)
			}
			return
		},
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
