package utils

import "log"

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
