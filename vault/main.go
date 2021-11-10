package vault

import (
	"log"
	"os"
)

var (
	DebugLog = log.New(os.Stdout, "debug#", log.Lshortfile)
	ErrorLog = log.New(os.Stderr, "error#", log.Lshortfile)
)
