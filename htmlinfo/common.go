package htmlinfo

import (
	"log"
	"os"
)

var (
	ErrorLogger = log.New(os.Stderr, `htmlinfo/error/`, log.Lshortfile)
	DebugLogger = log.New(os.Stdout, `htmlinfo/debug/`, log.Lshortfile)
	Debug = false
)