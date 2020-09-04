package main

import (
	"flag"
	"log"
	"os"
)

var (
	flagINPUT = flag.String(
		"i",
		"",
		"Input Mail:Pass file",
	)
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

	if *flagINPUT == "" {
		log.Fatalf("please define a input file\n%s -h for help\n", os.Args[0])
	}
}
