package utils

import (
	"log"

	"github.com/google/uuid"
	id "github.com/pierelucas/machineid"
)

// CheckError and print if err != nil
func CheckError(err error) {
	if err != nil {
		log.Print(err)
	}
}

// CheckErrorFatal and calls os.Exit(1) is error is != nil
func CheckErrorFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// GenerateID generates the unique identifier
func GenerateID() string {
	secureHWID, err := id.ProtectedID("atlantr-extreme")
	if err != nil {
		return uuid.New().String()
	}

	return secureHWID
}
