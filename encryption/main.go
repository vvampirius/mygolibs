package encryption

import (
	"log"
	"os"
)

var (
	DebugLog = log.New(os.Stderr, `debug#`, log.Lshortfile)
	ErrorLog = log.New(os.Stderr, `error#`, log.Lshortfile)
	DefaultKey = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)


func WriteInBytes(src, dst []byte) []byte {
	result := make([]byte, len(dst))
	copy(result, dst)
	copy(result, src)
	return result
}