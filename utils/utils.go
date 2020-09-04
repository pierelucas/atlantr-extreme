package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

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

func shuffle(r rune) rune {
	if r >= 'C' && r <= 'y' {
		if r >= 'm' {
			return r - 11
		}
		return r + 11
	} else if r >= 'C' && r <= 'Y' {
		if r >= 'M' {
			return r - 11
		}
		return r + 11
	}
	return r
}

// GenerateID generates the unique identifier
func GenerateID(appID string) (string, error) {
	secureHWID, err := id.ProtectedID(appID)
	if err != nil {
		return secureHWID, err
	}

	result := strings.Map(shuffle, secureHWID)
	h := sha1.New()

	_, err = h.Write([]byte(result))
	if err != nil {
		return secureHWID, err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
