package instagram

import (
        "log"
        "os"
)

var (
        DebugLog = log.New(os.Stderr, `debug#`, log.Lshortfile)
        ErrorLog = log.New(os.Stderr, `error#`, log.Lshortfile)
)
