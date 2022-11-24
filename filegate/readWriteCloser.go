package filegate

import (
	"os"
	"sync"
)

type ReadWriteCloser struct {
	fileFd *os.File
	fileGateUnlock func()
	mode uint8
	mu sync.Mutex
}

func (rwc *ReadWriteCloser) Read(p []byte) (int, error) {
	rwc.mu.Lock()
	defer rwc.mu.Unlock()
	if rwc.mode != 1 {
		rwc.fileFd.Seek(0, 0)
		rwc.mode = 1
	}
	return rwc.fileFd.Read(p)
}

func (rwc *ReadWriteCloser) Write(p []byte) (int, error) {
	rwc.mu.Lock()
	defer rwc.mu.Unlock()
	if rwc.mode != 2 {
		if err := rwc.fileFd.Truncate(0); err != nil {
			return 0, err
		}
		if _, err := rwc.fileFd.Seek(0, 0); err != nil {
			return 0, err
		}
		rwc.mode = 2
	}
	return rwc.fileFd.Write(p)
}

func (rwc *ReadWriteCloser) Close() error {
	rwc.mu.Lock()
	defer rwc.mu.Unlock()
	err := rwc.fileFd.Close()
	rwc.fileGateUnlock()
	return err
}

func (rwc *ReadWriteCloser) Reset() {
	rwc.mode = 0
}