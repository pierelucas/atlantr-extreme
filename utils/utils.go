package utils

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	id "github.com/pierelucas/machineid"
)

// Base64Encode --
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// CheckError and print to log if err != nil
func CheckError(err error) {
	if err != nil {
		log.Println(err)
	}
}

// CheckErrorFatal print to log and calls os.Exit(1) is error is != nil
func CheckErrorFatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// CheckErrorPrint and print if err != nil
func CheckErrorPrint(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// CheckErrorPrintFatal and calls os.Exit(1) is error is != nil
func CheckErrorPrintFatal(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// MultiLogf writes to standard logg and standard output
func MultiLogf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	log.Printf(format, a...)
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

// GotLineCount --
func GotLineCount(filepath string) (int32, error) {
	// We check if the file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return 0, err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var lineCount int32 = 0

	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		lineCount++
	}

	return lineCount, nil
}

// CheckDir checks if directory exists and creates it otherwise with perm 0755
func CheckDir(filepath string) error {
	var err error

	if _, err = os.Stat(filepath); os.IsNotExist(err) {
		err = os.Mkdir(filepath, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
