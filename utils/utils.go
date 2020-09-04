package utils

import (
	"encoding/base64"
	"log"

	"github.com/google/uuid"
	id "github.com/pierelucas/machineid"
)

// Base64Encode --
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

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
