package main

import "flag"

var (
	flagLOGOUTPUT = flag.String(
		"l",
		"",
		"Logoutput",
	)
	flagLASTLINELOG = flag.String(
		"ll",
		"",
		"lastlinelog",
	)
)

func init() {
	flag.Parse() // Parse our command-line arguments
}
